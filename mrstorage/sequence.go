package mrstorage

import "context"

type (
	// SequenceGenerator - генерирует последовательность из натуральных чисел.
	SequenceGenerator interface {
		Next(ctx context.Context) (nextID uint64, err error)
		MultiNext(ctx context.Context, count uint32) (nextIDs []uint64, err error)
	}
)
