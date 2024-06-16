package mrpostgres

import (
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLBuilderUpdate - comment struct.
	SQLBuilderUpdate struct {
		meta  *mrsql.EntityMetaUpdate
		set   *SQLBuilderSet
		where *SQLBuilderWhere
	}
)

// NewSQLBuilderUpdate - создаёт объект SQLBuilderUpdate.
func NewSQLBuilderUpdate(set *SQLBuilderSet, where *SQLBuilderWhere) *SQLBuilderUpdate {
	return &SQLBuilderUpdate{
		set:   set,
		where: where,
	}
}

// NewSQLBuilderUpdateWithMeta - создаёт объект SQLBuilderUpdate с метаинформацией.
func NewSQLBuilderUpdateWithMeta(meta *mrsql.EntityMetaUpdate, set *SQLBuilderSet, where *SQLBuilderWhere) *SQLBuilderUpdate {
	return &SQLBuilderUpdate{
		meta:  meta,
		set:   set,
		where: where,
	}
}

// Set - comment method.
func (b *SQLBuilderUpdate) Set(f func(s mrstorage.SQLBuilderSet) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.set != nil {
		partFunc = f(b.set)
	}

	return mrsql.NewBuilderPart(partFunc)
}

// SetFromEntity - comment method.
func (b *SQLBuilderUpdate) SetFromEntity(entity any) (mrstorage.SQLBuilderPart, error) {
	return b.SetFromEntityWith(entity, nil)
}

// SetFromEntityWith - comment method.
func (b *SQLBuilderUpdate) SetFromEntityWith(
	entity any,
	extFields func(s mrstorage.SQLBuilderSet) mrstorage.SQLBuilderPartFunc,
) (mrstorage.SQLBuilderPart, error) {
	if b.meta == nil {
		return nil, mrcore.ErrInternalNilPointer.New()
	}

	if b.set == nil {
		return nil, mrcore.ErrInternalNilPointer.New()
	}

	dbNames, args, err := b.meta.FieldsForUpdate(entity)
	if err != nil {
		return nil, err
	}

	if extFields == nil {
		return mrsql.NewBuilderPart(b.set.Fields(dbNames, args)), nil
	}

	return mrsql.NewBuilderPart(
		b.set.Join(
			b.set.Fields(dbNames, args),
			extFields(b.set),
		),
	), nil
}

// Where - comment method.
func (b *SQLBuilderUpdate) Where(f func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.where != nil {
		partFunc = f(b.where)
	} else {
		// protection against unintended updates
		partFunc = func(_ int) (string, []any) {
			return "1 = 0", nil
		}
	}

	return mrsql.NewBuilderPart(partFunc)
}
