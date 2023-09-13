package mrrabbitmq

import "github.com/mondegor/go-webcore/mrcore"

func (c *ConnAdapter) wrapError(err error) error {
    return mrcore.FactoryErrStorageQueryFailed.Caller(2).Wrap(err)
}
