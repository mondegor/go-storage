package mrfilestorage

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	FactoryErrInvalidPath = mrerr.NewFactoryWithCaller(
		"errMrFileStorageInvalidPath", mrerr.ErrorKindInternal, "invalid path '{{ .path }}'")
)
