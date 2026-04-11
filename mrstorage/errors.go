package mrstorage

import (
	"github.com/mondegor/go-sysmess/errors"
)

// ErrSystemFileProviderPingError - ошибка проверки работоспособности файлового провайдера.
var ErrSystemFileProviderPingError = errors.NewSystemProto("file provider ping error")
