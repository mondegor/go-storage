package mrpostgres

import (
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

type (
	SqlBuilderUpdate struct {
		meta  *mrsql.EntityMetaUpdate
		set   *SqlBuilderSet
		where *SqlBuilderWhere
	}
)

func NewSqlBuilderUpdate(set *SqlBuilderSet, where *SqlBuilderWhere) *SqlBuilderUpdate {
	return &SqlBuilderUpdate{
		set:   set,
		where: where,
	}
}

func NewSqlBuilderUpdateWithMeta(meta *mrsql.EntityMetaUpdate, set *SqlBuilderSet, where *SqlBuilderWhere) *SqlBuilderUpdate {
	return &SqlBuilderUpdate{
		meta:  meta,
		set:   set,
		where: where,
	}
}

func (b *SqlBuilderUpdate) Set(f func(s mrstorage.SqlBuilderSet) mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPart {
	var partFunc mrstorage.SqlBuilderPartFunc

	if b.set != nil {
		partFunc = f(b.set)
	}

	return mrsql.NewBuilderPart(partFunc)
}

func (b *SqlBuilderUpdate) SetFromEntity(entity any) (mrstorage.SqlBuilderPart, error) {
	return b.SetFromEntityWith(entity, nil)
}

func (b *SqlBuilderUpdate) SetFromEntityWith(entity any, extFields func(s mrstorage.SqlBuilderSet) mrstorage.SqlBuilderPartFunc) (mrstorage.SqlBuilderPart, error) {
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

func (b *SqlBuilderUpdate) Where(f func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPart {
	var partFunc mrstorage.SqlBuilderPartFunc

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
