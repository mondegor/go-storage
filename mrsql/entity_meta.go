package mrsql

import "github.com/mondegor/go-sysmess/mrlog"

type (
	// EntityMeta - объект для управления динамическим обновлением записей в БД.
	// Информация об обновлении считывается из тегов структуры.
	EntityMeta struct {
		metaUpdate  *EntityMetaUpdate
		metaOrderBy *EntityMetaOrderBy
	}
)

// ParseEntity - парсит указанную структуру entity и на основе её тегов
// создаёт объекты EntityMetaUpdate и EntityMetaOrderBy.
func ParseEntity(logger mrlog.Logger, entity any) (EntityMeta, error) {
	metaUpdate, err := NewEntityMetaUpdate(logger, entity)
	if err != nil {
		return EntityMeta{}, err
	}

	metaOrderBy, err := NewEntityMetaOrderBy(logger, entity)
	if err != nil {
		return EntityMeta{}, err
	}

	return EntityMeta{
		metaUpdate:  metaUpdate,
		metaOrderBy: metaOrderBy,
	}, nil
}

// MetaUpdate - возвращает метаинформацию об обновлении полей из распарсенной структуры.
func (e *EntityMeta) MetaUpdate() *EntityMetaUpdate {
	return e.metaUpdate
}

// MetaOrderBy - возвращает метаинформацию о сортировке полей из распарсенной структуры.
func (e *EntityMeta) MetaOrderBy() *EntityMetaOrderBy {
	return e.metaOrderBy
}
