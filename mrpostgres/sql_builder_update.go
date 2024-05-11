package mrpostgres

import (
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

type (
	SQLBuilderUpdate struct {
		meta  *mrsql.EntityMetaUpdate
		set   *SQLBuilderSet
		where *SQLBuilderWhere
	}
)

func NewSQLBuilderUpdate(set *SQLBuilderSet, where *SQLBuilderWhere) *SQLBuilderUpdate {
	return &SQLBuilderUpdate{
		set:   set,
		where: where,
	}
}

func NewSQLBuilderUpdateWithMeta(meta *mrsql.EntityMetaUpdate, set *SQLBuilderSet, where *SQLBuilderWhere) *SQLBuilderUpdate {
	return &SQLBuilderUpdate{
		meta:  meta,
		set:   set,
		where: where,
	}
}

func (b *SQLBuilderUpdate) Set(f func(s mrstorage.SQLBuilderSet) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.set != nil {
		partFunc = f(b.set)
	}

	return mrsql.NewBuilderPart(partFunc)
}

func (b *SQLBuilderUpdate) SetFromEntity(entity any) (mrstorage.SQLBuilderPart, error) {
	return b.SetFromEntityWith(entity, nil)
}

func (b *SQLBuilderUpdate) SetFromEntityWith(entity any, extFields func(s mrstorage.SQLBuilderSet) mrstorage.SQLBuilderPartFunc) (mrstorage.SQLBuilderPart, error) {
	if b.meta == nil {
		return nil, mrcore.FactoryErrInternalNilPointer.New()
	}

	if b.set == nil {
		return nil, mrcore.FactoryErrInternalNilPointer.New()
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

func (b *SQLBuilderUpdate) Where(f func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.where != nil {
		partFunc = f(b.where)
	} else {
		// protection against unintended updates
		partFunc = func(paramNumber int) (string, []any) {
			return "1 = 0", []any{}
		}
	}

	return mrsql.NewBuilderPart(partFunc)
}
