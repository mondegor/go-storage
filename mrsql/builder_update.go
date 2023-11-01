package mrsql

import (
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
)

type (
    BuilderUpdate struct {
        meta *EntityMetaUpdate
        set mrstorage.SqlBuilderSet
    }
)

func NewBuilderUpdate(set mrstorage.SqlBuilderSet) *BuilderUpdate {
    return &BuilderUpdate{
        set: set,
    }
}

func NewBuilderUpdateWithMeta(meta *EntityMetaUpdate, set mrstorage.SqlBuilderSet) *BuilderUpdate {
    return &BuilderUpdate{
        meta: meta,
        set: set,
    }
}

func (b *BuilderUpdate) Set(f func (s mrstorage.SqlBuilderSet) mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPart {
    return NewBuilderPart(f(b.set))
}

func (b *BuilderUpdate) SetFromEntity(entity any) (mrstorage.SqlBuilderPart, error) {
    if b.meta == nil {
        return nil, mrcore.FactoryErrInternalNilPointer.New()
    }

    dbNames, args, err := FieldsForUpdate(b.meta, entity)

    if err != nil {
        return nil, err
    }

    return NewBuilderPart(b.set.Fields(dbNames, args)), nil
}
