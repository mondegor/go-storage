package mrredis

import (
	"github.com/mondegor/go-webcore/mrcore"
    "github.com/redis/go-redis/v9"
)

func (c *Connection) wrapError(err error) error {
    if err == redis.Nil {
        return mrcore.FactoryErrNoRowFound.Caller(2).Wrap(err)
    }

    return mrcore.FactoryErrQueryFailed.Caller(2).Wrap(err)
}
