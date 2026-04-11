package mrstorage

import "context"

type (
	// SequenceGenerator - генератор последовательностей натуральных чисел.
	// Используется для получения уникальных идентификаторов (ID) из БД.
	SequenceGenerator interface {
		Next(ctx context.Context) (nextID uint64, err error)
		MultiNext(ctx context.Context, count int) (nextIDs []uint64, err error)
	}
)
