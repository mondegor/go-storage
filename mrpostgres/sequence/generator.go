package sequence

import (
	"context"
	"fmt"

	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-storage/mrstorage"
)

const (
	defaultMaxIDsOneQuery = 1024
)

type (
	// Generator - генератор последовательностей (на основе postgres).
	Generator struct {
		client                  mrstorage.DBConnManager
		maxIDsOneQuery          uint32
		sqlGeneratorSequenceID  string
		sqlGeneratorSequenceIDs string
	}
)

// NewGenerator - создаёт объект Generator.
func NewGenerator(client mrstorage.DBConnManager, sequenceName string, opts ...Option) *Generator {
	g := &Generator{
		client:                  client,
		maxIDsOneQuery:          defaultMaxIDsOneQuery,
		sqlGeneratorSequenceID:  `SELECT nextval('` + sequenceName + `');`,
		sqlGeneratorSequenceIDs: `SELECT nextval('` + sequenceName + `') FROM generate_series(1, $1);`,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Next - возвращает следующий свободный ID.
func (g Generator) Next(ctx context.Context) (nextID uint64, err error) {
	err = g.client.Conn(ctx).QueryRow(
		ctx,
		g.sqlGeneratorSequenceID,
	).Scan(
		&nextID,
	)

	return nextID, err
}

// MultiNext - возвращает нужное кол-во идентификаторов, но без гарантии непрерывности.
func (g Generator) MultiNext(ctx context.Context, count uint32) (nextIDs []uint64, err error) {
	if count < 2 {
		if count == 0 {
			return nil, mrcore.ErrInternalWithDetails.New("count must be greater than zero")
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

	for i := uint32(1); i <= batches; i++ {
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

	if count != uint32(len(nextIDs)) {
		return nil, mrcore.ErrStorageFetchDataFailed.Wrap(fmt.Errorf("expected next ids %d, got %d", count, len(nextIDs)))
	}

	return nextIDs, nil
}
