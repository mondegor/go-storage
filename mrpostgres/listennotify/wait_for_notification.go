package listennotify

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-sysmess/mrlog"
)

const (
	defaultReadyTimeout   = 5 * time.Second
	defaultReconnectDelay = 5 * time.Second
)

type (
	// ProcessWaitForNotification - процесс прослушивания и обработки событий (NOTIFY) от PostgreSQL.
	// Переподключается к БД при разрыве соединения с настраиваемой задержкой.
	ProcessWaitForNotification struct {
		conn               *mrpostgres.ConnAdapter
		logger             mrlog.Logger
		listenerChannelMap map[string]chan struct{} // listenerChannelMap - маппинг имён каналов на каналы уведомлений
		reconnectDelay     time.Duration            // reconnectDelay - задержка между попытками переподключения

		wg   sync.WaitGroup
		done chan struct{}

		receiverChannels receiveChannels // receiverChannels - коллекция каналов для подписчиков
	}
)

// NewProcessWaitForNotification - создаёт объект ProcessWaitForNotification для прослушивания NOTIFY от PostgreSQL.
// Параметры:
//   - conn - адаптер подключения к PostgreSQL;
//   - logger - логгер для вывода сообщений;
//   - channels - список имён каналов для подписки.
func NewProcessWaitForNotification(
	conn *mrpostgres.ConnAdapter,
	logger mrlog.Logger,
	channels []string,
) *ProcessWaitForNotification {
	listenerChannelMap, receiverChannelList := createListenerChannels(logger, channels)

	return &ProcessWaitForNotification{
		conn:               conn,
		logger:             logger,
		listenerChannelMap: listenerChannelMap,
		reconnectDelay:     defaultReconnectDelay,

		wg:   sync.WaitGroup{},
		done: make(chan struct{}),

		receiverChannels: receiverChannelList,
	}
}

// Caption - возвращает название процесса в свободной форме.
func (p *ProcessWaitForNotification) Caption() string {
	return "ProcessWaitForNotification"
}

// ReadyTimeout - возвращает таймаут готовности процесса для ожидания запуска.
func (p *ProcessWaitForNotification) ReadyTimeout() time.Duration {
	return defaultReadyTimeout
}

// Start - запускает процесс прослушивания NOTIFY от PostgreSQL.
// Блокирует выполнение до завершения контекста или возникновения ошибки.
// Автоматически переподключается при разрыве соединения с задержкой defaultReconnectDelay.
func (p *ProcessWaitForNotification) Start(ctx context.Context, ready func()) error {
	p.wg.Add(1)
	defer p.wg.Done()

	p.logger.Debug(ctx, "Starting the WaitForNotification...")
	defer p.logger.Debug(ctx, "The WaitForNotification has been stopped")

	ctxListen, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case <-p.done:
			cancel()
		case <-ctx.Done():
		}
	}()

	if ready != nil {
		ready()
	}

	for {
		if err := p.listen(ctxListen); err != nil {
			if errors.Is(err, ctxListen.Err()) {
				return nil
			}

			p.logger.Error(ctxListen, "ProcessWaitForNotification.listen", "error", err)
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
			p.logger.Debug(ctx, "The WaitForNotification detected context 'Done'", "error", ctx.Err())

			return nil
		case <-time.After(p.reconnectDelay):
		}
	}
}

// Find - находит канал по имени и возвращает его для получения уведомлений.
func (p *ProcessWaitForNotification) Find(name string) (<-chan struct{}, error) {
	return p.receiverChannels.Find(name)
}

// MustFind - находит канал по имени и возвращает его для получения уведомлений.
// Если имя канала не зарегистрировано, то регистрируется ошибка и возвращается
// фиктивный канал, который будет заблокирован до завершения процесса.
func (p *ProcessWaitForNotification) MustFind(name string) <-chan struct{} {
	ch, err := p.receiverChannels.Find(name)
	if err != nil {
		mrlog.Error(p.logger, "ProcessWaitForNotification.MustFind", "error", err)

		return p.done
	}

	return ch
}

// Shutdown - корректно завершает процесс прослушивания NOTIFY.
func (p *ProcessWaitForNotification) Shutdown(ctx context.Context) error {
	p.logger.Info(ctx, "Shutting down the WaitForNotification...")
	close(p.done)

	p.wg.Wait()
	p.logger.Info(ctx, "The WaitForNotification has been shut down")

	return nil
}

func (p *ProcessWaitForNotification) listen(ctx context.Context) error {
	conn, err := p.conn.HijackConn(ctx)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	defer func() {
		_ = conn.Close(ctx)
	}()

	for name := range p.listenerChannelMap {
		if _, err := conn.Exec(ctx, "LISTEN "+pgx.Identifier{name}.Sanitize()); err != nil {
			return fmt.Errorf("unable to start listening channel '%s': %w", name, err)
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
				p.logger.Debug(ctx, fmt.Sprintf("Received notification: PID=%d, Channel='%s', Payload='%s'", note.PID, note.Channel, note.Payload))
			default:
				p.logger.Info(ctx, fmt.Sprintf("Double notification: PID=%d, Channel='%s', Payload='%s' [skipped]", note.PID, note.Channel, note.Payload))
			}
		} else {
			p.logger.Warn(
				ctx,
				"Unknown channel",
				"pid", note.PID,
				"channel", note.Channel,
				"payload", note.Payload,
			)
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

func createListenerChannels(logger mrlog.Logger, channels []string) (map[string]chan struct{}, []receiveChannel) {
	listenerChannels := make(map[string]chan struct{}, len(channels))
	receiveChannels := make([]receiveChannel, 0, len(channels))

	for _, name := range channels {
		if _, ok := listenerChannels[name]; ok {
			mrlog.Warn(logger, "Duplicate listen channel", "channel", name)

			continue
		}

		channel := make(chan struct{})

		listenerChannels[name] = channel
		receiveChannels = append(
			receiveChannels,
			receiveChannel{
				Name:    name,
				Channel: channel,
			},
		)
	}

	return listenerChannels, receiveChannels
}
