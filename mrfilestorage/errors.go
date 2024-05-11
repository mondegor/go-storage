package mrfilestorage

import e "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrInvalidPath = e.NewFactory(
		"errMrFileStorageInvalidPath", e.ErrorTypeInternal, "invalid path '{{ .path }}'")
)
