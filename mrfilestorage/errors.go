package mrfilestorage

import "github.com/mondegor/go-sysmess/mrerr"

// ErrInvalidPath - invalid path.
var ErrInvalidPath = mrerr.NewProto(
	"errMrFileStorageInvalidPath", mrerr.ErrorKindInternal, "invalid path '{{ .path }}'")
