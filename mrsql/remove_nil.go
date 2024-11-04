package mrsql

import (
	"github.com/mondegor/go-storage/mrstorage"
)

// SQLPartFuncRemoveNil - уменьшает указанный массив удаляя из него все nil элементы.
func SQLPartFuncRemoveNil(parts []mrstorage.SQLPartFunc) []mrstorage.SQLPartFunc {
	needOffset := false
	length := 0

	for i := range parts {
		if parts[i] == nil {
			needOffset = true

			continue
		}

		if i > length {
			parts[length] = parts[i]
		}

		length++
	}

	if needOffset {
		parts = parts[:length]
	}

	return parts
}
