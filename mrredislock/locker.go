package mrredislock

import (
    "context"
    "time"

    "github.com/bsm/redislock"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/redis/go-redis/v9"
)

// go get -u github.com/bsm/redislock

type (
    lockerAdapter struct {
        lock *redislock.Client
    }
)

func NewLockerAdapter(conn redis.UniversalClient) *lockerAdapter {
    return &lockerAdapter{
        lock: redislock.New(conn),
    }
}

func (l *lockerAdapter) Lock(ctx context.Context, key string) (mrcore.UnlockFunc, error) {
    return l.LockWithExpiry(ctx, key, 0)
}

// LockWithExpiry - if expiry = 0 then set expiry by default
func (l *lockerAdapter) LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (mrcore.UnlockFunc, error) {
    if expiry == 0 {
        expiry = mrcore.LockerDefaultExpiry
    }

    mutex, err := l.lock.Obtain(ctx, key, expiry, nil)

    if err != nil {
        return nil, mrcore.FactoryErrInternal.Wrap(err)
    }

    return func() {
        err := mutex.Release(ctx)

        if err != nil {
            mrctx.Logger(ctx).Error(
                "mrredislock.lockerAdapter::MutexUnlock=%s; err: %s",
                key,
                err,
            )
        }
    }, nil
}