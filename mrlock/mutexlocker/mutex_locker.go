package mutexlocker

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/mondegor/go-storage/mrlock"
)

const (
	lockerName = "MutexLocker"
)

type (
	// Locker - реализует интерфейс блокировщика указанного ключа основанный на mutex.
	Locker struct {
		logger logger
		tracer tracer
		mu     sync.Mutex
		keys   map[string]int64
	}

	logger interface {
		Warn(ctx context.Context, msg string, args ...any)
	}

	tracer interface {
		Trace(ctx context.Context, args ...any)
	}
)

// New - создаёт объект Locker.
func New(logger logger, tracer tracer, minBufferSize int) *Locker {
	return &Locker{
		tracer: tracer,
		logger: logger,
		keys:   make(map[string]int64, minBufferSize),
	}
}

// Lock - блокирует указанный ключ в рамках приложения с временем блокировки по умолчанию
// и возвращает функцию разблокировки.
func (l *Locker) Lock(ctx context.Context, key string) (unlock func(), err error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - блокирует указанный ключ в рамках приложения с указанием
// времени блокировки и возвращает функцию разблокировки.
// Если указана expiry равная нулю, то будет установлено время блокировки по умолчанию.
func (l *Locker) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (unlock func(), err error) {
	if expiry == 0 {
		expiry = mrlock.DefaultExpiry
	}

	l.traceCmd(ctx, "Lock:"+expiry.String()+", keys-len="+strconv.Itoa(len(l.keys)), key)

	l.mu.Lock()
	defer l.mu.Unlock()

	if exp, ok := l.keys[key]; ok && exp > time.Now().UnixNano() {
		return nil, mrlock.ErrStorageLockKeyNotObtained.New(
			"source", lockerName,
			"lock_key", key,
		)
	}

	l.keys[key] = time.Now().UnixNano() + expiry.Nanoseconds()

	return func() {
		l.traceCmd(ctx, "Unlock", key)

		l.mu.Lock()
		defer l.mu.Unlock()

		if _, ok := l.keys[key]; ok {
			delete(l.keys, key)
		} else {
			l.logger.Warn(
				ctx,
				"unlock",
				"error", mrlock.ErrStorageLockKeyNotHeld.New(
					"source", lockerName,
					"lock_key", key,
				),
			)
		}
	}, nil
}

func (l *Locker) traceCmd(ctx context.Context, command, key string) {
	l.tracer.Trace(
		ctx,
		"source", lockerName,
		"cmd", command,
		"key", key,
	)
}
