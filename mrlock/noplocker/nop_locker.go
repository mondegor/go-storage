package noplocker

import (
	"context"
	"time"

	"github.com/mondegor/go-storage/mrlock"
)

const (
	// lockerName - имя блокировщика для логирования и трассировки.
	lockerName = "NopLocker"
)

type (
	// Locker - заглушка, реализующая интерфейс блокировщика ключа.
	// Не выполняет реальной блокировки, только логирует операции через tracer.
	// Полезен для тестирования и в средах, где блокировка не требуется.
	Locker struct {
		tracer tracer
	}

	tracer interface {
		Trace(ctx context.Context, args ...any)
	}
)

// New - создаёт объект Locker-заглушку, не выполняющую реальной блокировки.
func New(tracer tracer) *Locker {
	return &Locker{
		tracer: tracer,
	}
}

// Lock - эмулирует блокировку указанного ключа с временем блокировки по умолчанию
// и возвращает функцию разблокировки.
func (l *Locker) Lock(ctx context.Context, key string) (unlock func(), err error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - эмулирует блокировку указанного ключа
// с указанием времени блокировки и возвращает функцию разблокировки.
func (l *Locker) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (unlock func(), err error) {
	if expiry == 0 {
		expiry = mrlock.DefaultExpiry
	}

	l.traceCmd(ctx, "Lock:"+expiry.String(), key)

	return func() {
		l.traceCmd(ctx, "Unlock", key)
	}, nil
}

// traceCmd - логирует выполняемую операцию блокировки для трассировки.
func (l *Locker) traceCmd(ctx context.Context, command, key string) {
	l.tracer.Trace(
		ctx,
		"source", lockerName,
		"cmd", command,
		"key", key,
	)
}
