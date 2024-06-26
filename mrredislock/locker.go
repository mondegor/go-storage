package mrredislock

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/mondegor/go-webcore/mrlock"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/redis/go-redis/v9"
)

// https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html
// http://antirez.com/news/101

// go get -u github.com/bsm/redislock

const (
	lockerName = "RedisLock"
)

type (
	// LockerAdapter - адаптер для работы с сетевыми блокировками на основе Redis.
	LockerAdapter struct {
		lock *redislock.Client
	}
)

// NewLockerAdapter - создаёт объект LockerAdapter.
func NewLockerAdapter(conn redis.UniversalClient) *LockerAdapter {
	return &LockerAdapter{
		lock: redislock.New(conn),
	}
}

// Lock - comment method.
func (l *LockerAdapter) Lock(ctx context.Context, key string) (func(), error) {
	return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - if expiry = 0 then set expiry by default.
func (l *LockerAdapter) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (func(), error) {
	if expiry == 0 {
		expiry = mrlock.DefaultExpiry
	}

	l.traceCmd(ctx, "Lock:"+expiry.String(), key)

	mutex, err := l.lock.Obtain(ctx, key, expiry, nil)
	if err != nil {
		return nil, l.wrapError(err)
	}

	return func() {
		l.traceCmd(ctx, "Unlock", key)

		if err := mutex.Release(ctx); err != nil {
			mrlog.Ctx(ctx).
				Debug().
				Str("source", lockerName).
				Str("cmd", "unlock").
				Str("key", key).
				Err(err).
				Send()
		}
	}, nil
}
