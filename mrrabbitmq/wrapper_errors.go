package mrrabbitmq

import "github.com/mondegor/go-webcore/mrcore"

func (c *Connection) wrapError(err error) error {
    return mrcore.FactoryErrStorageQueryFailed.Caller(2).Wrap(err)
}
