package mrfilestorage

import "github.com/mondegor/go-sysmess/mrerr"

// ErrInvalidPath - comment var.
var ErrInvalidPath = mrerr.NewProto(
	"errMrFileStorageInvalidPath", mrerr.ErrorKindInternal, "invalid path '{{ .path }}'")
