package mrrabbitmq

import "github.com/mondegor/go-webcore/mrcore"

func (c *Connection) wrapError(err error) error {
    return mrcore.FactoryErrQueryFailed.Caller(2).Wrap(err)
}
