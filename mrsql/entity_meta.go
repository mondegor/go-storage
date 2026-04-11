package mrsql

import "github.com/mondegor/go-sysmess/mrlog"

type (
	// EntityMeta - метаинформация о структуре для динамического управления записями в БД.
	// Объединяет информацию об обновлении полей (EntityMetaUpdate) и сортировке (EntityMetaOrderBy).
	// Информация считывается из тегов структуры: `db`, `upd`, `sort`.
	EntityMeta struct {
		metaUpdate  *EntityMetaUpdate
		metaOrderBy *EntityMetaOrderBy
	}
)

// ParseEntity - парсит указанную структуру entity и на основе её тегов
// создаёт объекты EntityMetaUpdate и EntityMetaOrderBy.
// Параметры:
//   - logger - логгер для вывода информации о разборе тегов;
//   - entity - структура для разбора (может быть указателем).
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
// Используется для формирования части SET в SQL-запросе UPDATE.
func (e *EntityMeta) MetaUpdate() *EntityMetaUpdate {
	return e.metaUpdate
}

// MetaOrderBy - возвращает метаинформацию о сортировке полей из распарсенной структуры.
// Используется для формирования части ORDER BY в SQL-запросе.
func (e *EntityMeta) MetaOrderBy() *EntityMetaOrderBy {
	return e.metaOrderBy
}
