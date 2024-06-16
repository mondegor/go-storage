package mrpostgres

import (
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLBuilderCondition - comment struct.
	SQLBuilderCondition struct {
		where *SQLBuilderWhere
	}
)

// NewSQLBuilderCondition - создаёт объект SQLBuilderCondition.
func NewSQLBuilderCondition(where *SQLBuilderWhere) *SQLBuilderCondition {
	return &SQLBuilderCondition{
		where: where,
	}
}

// Where - comment method.
func (b *SQLBuilderCondition) Where(f func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.where != nil {
		partFunc = f(b.where)
	}

	return mrsql.NewBuilderPart(partFunc)
}
