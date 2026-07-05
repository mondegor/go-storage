package mrfilestorage

import (
	"github.com/mondegor/go-core/errors"
)

// ErrInternalInvalidPath - недопустимый путь к файлу.
var ErrInternalInvalidPath = errors.NewInternalProto("invalid path")
