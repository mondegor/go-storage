package mrredsync

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/redis/go-redis/v9"
)

// go get -u github.com/go-redsync/redsync/v4

const (
	lockerName = "RedSync"
)

type (
	lockerAdapter struct {
		lock *redsync.Redsync
	}
)

func NewLockerAdapter(conn redis.UniversalClient) *lockerAdapter {
	pool := goredis.NewPool(conn)

	return &lockerAdapter{
		lock: redsync.New(pool),
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

	l.debugCmd(ctx, "Lock:"+expiry.String(), key)

	options := []redsync.Option{
		redsync.WithExpiry(expiry),
	}

	mutex := l.lock.NewMutex(key, options...)

	if err := mutex.LockContext(ctx); err != nil {
		return nil, l.wrapError(err)
	}

	return func() {
		l.debugCmd(ctx, "Unlock", key)

		if _, err := mutex.UnlockContext(ctx); err != nil {
			mrctx.Logger(ctx).Error(
				"%s: cmd=unlock, key=%s, err={%s}",
				lockerName,
				key,
				err,
			)
		}
	}, nil
}
