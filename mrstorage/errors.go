package mrstorage

import "github.com/mondegor/go-sysmess/mrerr"

// ErrFileProviderPingError - file provider ping error.
var ErrFileProviderPingError = mrerr.NewKindSystem("file provider ping error: '{Name}'")
