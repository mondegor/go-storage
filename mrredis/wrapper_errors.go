package mrredis

import (
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/redis/go-redis/v9"
)

func (c *Connection) wrapError(err error) error {
    if err == redis.Nil {
        return mrstorage.ErrFactoryNoRowFound.Caller(2).Wrap(err)
    }

    return mrstorage.ErrFactoryQueryFailed.Caller(2).Wrap(err)
}
