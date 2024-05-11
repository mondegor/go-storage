package mrrabbitmq

import "github.com/mondegor/go-webcore/mrcore"

func (c *ConnAdapter) wrapError(err error) error {
	const skipFrame = 1
	return mrcore.FactoryErrStorageQueryFailed.WithSkipFrame(skipFrame).Wrap(err)
}
