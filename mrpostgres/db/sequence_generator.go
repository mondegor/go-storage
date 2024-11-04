package db

import (
	"context"
	"errors"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SequenceGenerator - генератор как одного ID, так и указанного кол-ва (на основе postgres последовательностей).
	SequenceGenerator struct {
		client                  mrstorage.DBConnManager
		sqlGeneratorSequenceID  string
		sqlGeneratorSequenceIDs string
	}
)

// NewSequenceGenerator - создаёт объект SequenceGenerator.
func NewSequenceGenerator(client mrstorage.DBConnManager, sequenceName string) SequenceGenerator {
	return SequenceGenerator{
		client:                  client,
		sqlGeneratorSequenceID:  `SELECT nextval('` + sequenceName + `');`,
		sqlGeneratorSequenceIDs: `SELECT setval('` + sequenceName + `', nextval('` + sequenceName + `') + $1);`,
	}
}

// Next - возвращает следующий свободный ID.
func (r SequenceGenerator) Next(ctx context.Context) (nextID uint64, err error) {
	err = r.client.Conn(ctx).QueryRow(
		ctx,
		r.sqlGeneratorSequenceID,
	).Scan(
		&nextID,
	)

	return nextID, err
}

// Reserve - резервирует указанное кол-во ID и возвращает первый зарезервированный ID.
func (r SequenceGenerator) Reserve(ctx context.Context, count uint16) (firstID uint64, err error) {
	if count == 0 {
		return 0, errors.New("count must be greater than zero")
	}

	err = r.client.Conn(ctx).QueryRow(
		ctx,
		r.sqlGeneratorSequenceIDs,
		count-1,
	).Scan(
		&firstID,
	)

	return firstID, err
}
