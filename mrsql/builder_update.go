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
    return b.SetFromEntityWith(entity, nil)
}

func (b *BuilderUpdate) SetFromEntityWith(entity any, extFields func(s mrstorage.SqlBuilderSet) mrstorage.SqlBuilderPartFunc) (mrstorage.SqlBuilderPart, error) {
    if b.meta == nil {
        return nil, mrcore.FactoryErrInternalNilPointer.New()
    }

    dbNames, args, err := b.meta.FieldsForUpdate(entity)

    if err != nil {
        return nil, err
    }

    if extFields == nil {
        return NewBuilderPart(b.set.Fields(dbNames, args)), nil
    }

    return NewBuilderPart(
        b.set.Join(
            b.set.Fields(dbNames, args),
            extFields(b.set),
        ),
    ), nil
}

