package locker

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/redis/go-redis/v9"
)

// https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html
// http://antirez.com/news/101

// go get -u github.com/bsm/redislock

const (
	// lockerName - имя блокировщика для логирования и трассировки.
	lockerName = "RedisLock"
)

type (
	// Adapter - адаптер для работы с сетевыми блокировками на основе Redis.
	Adapter struct {
		lock   *redislock.Client
		logger mrlog.Logger
		tracer mrtrace.Tracer
	}
)

// NewAdapter - создаёт объект Adapter.
func NewAdapter(
	conn redis.UniversalClient,
	logger mrlog.Logger,
	tracer mrtrace.Tracer,
) *Adapter {
	return &Adapter{
		lock:   redislock.New(conn),
		logger: logger,
		tracer: tracer,
	}
}

// Lock - захватывает блокировку указанного ключа с временем жизни по умолчанию.
// Возвращает функцию для освобождения блокировки.
func (l *Adapter) Lock(ctx context.Context, key string) (func(), error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - захватывает блокировку указанного ключа с заданным временем жизни.
// Возвращает функцию для освобождения блокировки.
func (l *Adapter) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (func(), error) {
	if expiry == 0 {
		expiry = mrlock.DefaultExpiry
	}

	l.traceCmd(ctx, "Lock:"+expiry.String(), key)

	mutex, err := l.lock.Obtain(ctx, key, expiry, nil)
	if err != nil {
		return nil, l.wrapError(err, key)
	}

	return func() {
		l.traceCmd(ctx, "Unlock", key)

		if err := mutex.Release(ctx); err != nil {
			l.logger.Warn(ctx, "unlock", "error", l.wrapError(err, key))
		}
	}, nil
}
