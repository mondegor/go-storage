package mrrabbitmq

import (
    "github.com/mondegor/go-storage/mrstorage"
)

func (c *Connection) wrapError(err error) error {
    return mrstorage.FactoryQueryFailed.Caller(2).Wrap(err)
}
