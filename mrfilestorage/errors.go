package mrfilestorage

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	FactoryErrInvalidPath = mrerr.NewFactory(
		"errMrFileStorageInvalidPath", mrerr.ErrorKindInternal, "invalid path '{{ .path }}'")
)
