package sequence

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"

	"github.com/mondegor/go-storage/mrstorage"
)

const (
	// defaultMaxIDsOneQuery - максимальное количество ID, получаемых за один запрос по умолчанию.
	defaultMaxIDsOneQuery = 1000
)

type (
	// Generator - генератор последовательностей на основе PostgreSQL.
	// Использует функции nextval() для получения уникальных ID из последовательностей БД.
	Generator struct {
		client                  mrstorage.DBConnManager // client - менеджер подключений к БД
		maxIDsOneQuery          int                     // maxIDsOneQuery - максимальное количество ID за один запрос
		sqlGeneratorSequenceID  string                  // sqlGeneratorSequenceID - SQL-запрос для получения одного ID
		sqlGeneratorSequenceIDs string                  // sqlGeneratorSequenceIDs - SQL-запрос для получения нескольких ID
	}
)

// NewGenerator - создаёт объект Generator для генерации уникальных ID из PostgreSQL-последовательности.
func NewGenerator(client mrstorage.DBConnManager, sequenceName string, opts ...Option) *Generator {
	o := options{
		generator: &Generator{
			client:                  client,
			sqlGeneratorSequenceID:  `SELECT nextval('` + sequenceName + `');`,
			sqlGeneratorSequenceIDs: `SELECT nextval('` + sequenceName + `') FROM generate_series(1, $1);`,
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.generator.maxIDsOneQuery < 1 {
		o.generator.maxIDsOneQuery = defaultMaxIDsOneQuery
	}

	return o.generator
}

// Next - возвращает следующий свободный ID из последовательности PostgreSQL.
// Использует функцию nextval() для атомарного получения ID.
func (g Generator) Next(ctx context.Context) (nextID uint64, err error) {
	err = g.client.Conn(ctx).QueryRow(
		ctx,
		g.sqlGeneratorSequenceID,
	).Scan(
		&nextID,
	)

	return nextID, err
}

// MultiNext - возвращает указанное количество идентификаторов из последовательности.
// ID получаются пакетами (batch) для оптимизации запросов к БД.
// Не гарантирует непрерывность последовательности (могут быть пропуски).
func (g Generator) MultiNext(ctx context.Context, count int) (nextIDs []uint64, err error) {
	if count < 2 {
		if count == 0 {
			return nil, errors.NewInternalError("count must be greater than zero")
		}

		nextID, err := g.Next(ctx)
		if err != nil {
			return nil, err
		}

		return []uint64{nextID}, nil
	}

	nextIDs = make([]uint64, 0, count)

	idsOneQuery := g.maxIDsOneQuery
	batches := count / idsOneQuery // кол-во полных запросов
	rest := count % idsOneQuery    // последний запрос

	if rest > 0 {
		batches++
	}

	for i := 1; i <= batches; i++ {
		if i == batches && rest > 0 {
			idsOneQuery = rest
		}

		err = func() error {
			cursor, err := g.client.Conn(ctx).Query(
				ctx,
				g.sqlGeneratorSequenceIDs,
				idsOneQuery,
			)
			if err != nil {
				return err
			}

			defer cursor.Close()

			for cursor.Next() {
				var nextID uint64

				err = cursor.Scan(
					&nextID,
				)
				if err != nil {
					return err
				}

				nextIDs = append(nextIDs, nextID)
			}

			return cursor.Err()
		}()
		if err != nil {
			return nil, err
		}
	}

	if count != len(nextIDs) {
		return nil, errors.ErrInternalStorageFetchDataFailed.WithDetails(
			"count != len(nextIDs)",
			"count", count,
			"actual", len(nextIDs),
			"source", "mrpostgres.SequenceGenerator",
		)
	}

	return nextIDs, nil
}
