package mrpostgres

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mondegor/go-sysmess/mrlog"
)

const (
	reconnectDelay = 5 * time.Second
)

type (
	// ReceiverChannel - comment struct.
	ReceiverChannel struct {
		Name    string
		Channel <-chan struct{}
	}

	// ReceiverChannels - comment array.
	ReceiverChannels []ReceiverChannel
)

// Find - comment method.
func (rc *ReceiverChannels) Find(name string) (<-chan struct{}, error) {
	for _, rch := range *rc {
		if name == rch.Name {
			return rch.Channel, nil
		}
	}

	return nil, fmt.Errorf("no such channel with name '%s'", name)
}

// MustFind - comment method.
func (rc *ReceiverChannels) MustFind(name string) <-chan struct{} {
	ch, err := rc.Find(name)
	if err != nil {
		panic(err)
	}

	return ch
}

type (
	// ProcessWaitForNotification - объект для прослушивания и обработки событий присылаемых БД.
	ProcessWaitForNotification struct {
		conn               *ConnAdapter
		logger             mrlog.Logger
		listenerChannelMap map[string]chan struct{}
		reconnectDelay     time.Duration

		wgMain sync.WaitGroup
		done   chan struct{}

		ReceiverChannels ReceiverChannels
	}
)

// NewProcessWaitForNotification - создаёт объект ProcessWaitForNotification.
func NewProcessWaitForNotification(
	conn *ConnAdapter,
	logger mrlog.Logger,
	channels []string,
) *ProcessWaitForNotification {
	listenerChannelMap, receiverChannels := createListenerChannels(channels)

	return &ProcessWaitForNotification{
		conn:               conn,
		logger:             logger,
		listenerChannelMap: listenerChannelMap,
		reconnectDelay:     reconnectDelay,

		wgMain: sync.WaitGroup{},
		done:   make(chan struct{}),

		ReceiverChannels: receiverChannels,
	}
}

// Caption - comment struct.
func (p *ProcessWaitForNotification) Caption() string {
	return "ProcessWaitForNotification"
}

// ReadyTimeout - comment struct.
func (p *ProcessWaitForNotification) ReadyTimeout() time.Duration {
	return 5 * time.Second
}

// Start - comment struct.
func (p *ProcessWaitForNotification) Start(ctx context.Context, ready func()) error {
	p.wgMain.Add(1)
	defer p.wgMain.Done()

	p.logger.Debug(ctx, "Starting the WaitForNotification...")
	defer p.logger.Debug(ctx, "The WaitForNotification has been stopped")

	if ready != nil {
		ready()
	}

	for {
		if err := p.listen(ctx); err != nil {
			if errors.Is(err, ctx.Err()) {
				return nil
			}

			p.logger.Error(ctx, "ProcessWaitForNotification.listen", "error", err)
		} else {
			return nil
		}

		if p.reconnectDelay < 1 {
			continue
		}

		select {
		case <-p.done:
			return nil
		case <-ctx.Done():
			return nil
		case <-time.After(p.reconnectDelay):
		}
	}
}

// Shutdown - comment struct.
func (p *ProcessWaitForNotification) Shutdown(ctx context.Context) error {
	p.logger.Info(ctx, "Shutting down the WaitForNotification...")
	close(p.done)

	p.wgMain.Wait()
	p.logger.Info(ctx, "The WaitForNotification has been shut down")

	return nil
}

func (p *ProcessWaitForNotification) listen(ctx context.Context) error {
	conn, err := p.conn.HijackConn(ctx)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close(ctx)

	for name := range p.listenerChannelMap {
		if _, err := conn.Exec(ctx, "LISTEN "+pgx.Identifier{name}.Sanitize()); err != nil {
			return fmt.Errorf("unable to start listening channel %s: %w", name, err)
		}
	}

	for {
		note, err := conn.WaitForNotification(ctx)
		if err != nil {
			return fmt.Errorf("waiting for notification: %w", err)
		}

		if ch, ok := p.listenerChannelMap[note.Channel]; ok {
			// если канал занят, значит такое же событие ещё не обработано,
			// поэтому нет смысла отправлять повторное событие
			select {
			case ch <- struct{}{}:
				p.logger.Debug(ctx, fmt.Sprintf("Received notification: PID=%d, Channel=%s, Payload=%s", note.PID, note.Channel, note.Payload))
			default:
				p.logger.Info(ctx, fmt.Sprintf("Double notification: PID=%d, Channel=%s, Payload=%s [skipped]", note.PID, note.Channel, note.Payload))
			}
		} else {
			p.logger.Warn(ctx, fmt.Sprintf("Unknown channel: PID=%d, Channel=%s, Payload=%s", note.PID, note.Channel, note.Payload))
		}

		select {
		case <-p.done:
			return nil
		default:
		}
	}
}

func createListenerChannels(channels []string) (map[string]chan struct{}, []ReceiverChannel) {
	listenerChannels := make(map[string]chan struct{}, len(channels))
	receiveChannels := make([]ReceiverChannel, 0, len(channels))

	for _, name := range channels {
		if _, ok := listenerChannels[name]; ok {
			// TODO: можно логировать
			continue
		}

		channel := make(chan struct{})

		listenerChannels[name] = channel
		receiveChannels = append(
			receiveChannels,
			ReceiverChannel{
				Name:    name,
				Channel: channel,
			},
		)
	}

	return listenerChannels, receiveChannels
}
