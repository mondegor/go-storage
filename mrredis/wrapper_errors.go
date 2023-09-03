package mrredis

import (
    "github.com/mondegor/go-storage/mrstorage"
    redislib "github.com/redis/go-redis/v9"
)

func (c *Connection) wrapError(err error) error {
    if err == redislib.Nil {
        return mrstorage.FactoryNoRowFound.Caller(2).Wrap(err)
    }

    return mrstorage.FactoryQueryFailed.Caller(2).Wrap(err)
}
