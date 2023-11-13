package mrredis

import (
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/redis/go-redis/v9"
)

func (c *ConnAdapter) wrapError(err error) error {
	if err == redis.Nil {
		return mrcore.FactoryErrStorageNoRowFound.Caller(1).Wrap(err)
	}

	return mrcore.FactoryErrStorageQueryFailed.Caller(1).Wrap(err)
}
