package mrfilestorage

import (
	"github.com/mondegor/go-sysmess/errors"
)

// ErrInternalInvalidPath - недопустимый путь к файлу.
var ErrInternalInvalidPath = errors.NewInternalProto("invalid path")
