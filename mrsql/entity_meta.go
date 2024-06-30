package mrsql

import (
	"regexp"

	"github.com/mondegor/go-storage/mrstorage"
)

var regexpDbName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

type (
	// EntityMeta - метаинформация о таблице БД, которую используют сторонние модули.
	EntityMeta struct {
		tableName   string
		primaryName string
		where       mrstorage.SQLBuilderPart
	}
)

// NewEntityMeta - создаёт объект EntityMeta.
func NewEntityMeta(tableName, primaryName string, where mrstorage.SQLBuilderPart) *EntityMeta {
	return &EntityMeta{
		tableName:   tableName,
		primaryName: primaryName,
		where:       where,
	}
}

// TableName - comment method.
func (e *EntityMeta) TableName() string {
	return e.tableName
}

// PrimaryName - comment method.
func (e *EntityMeta) PrimaryName() string {
	return e.primaryName
}

// Condition - comment method.
func (e *EntityMeta) Condition() mrstorage.SQLBuilderPart {
	return e.where
}
