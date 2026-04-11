package mrsql

import (
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// Part - часть SQL-запроса с используемыми в ней аргументами.
	// Позволяет формировать SQL с настраиваемым префиксом и номером начального аргумента.
	Part struct {
		argumentNumber int                   // argumentNumber - номер первого аргумента для нумерации ($1, $2, ...)
		sqlPrefix      string                // sqlPrefix - префикс, добавляемый перед SQL-выражением (например: " WHERE ")
		partFunc       mrstorage.SQLPartFunc // partFunc - функция для генерации SQL-выражения
	}
)

// NewPart - создаёт объект Part с указанным номером начального аргумента и функцией генерации SQL.
// Параметры:
//   - argumentNumber - номер первого аргумента для нумерации параметров ($1, $2, ...);
//   - part - функция, возвращающая SQL-выражение и список аргументов.
func NewPart(argumentNumber int, part mrstorage.SQLPartFunc) *Part {
	return &Part{
		argumentNumber: argumentNumber,
		partFunc:       part,
	}
}

// WithPrefix - возвращает копию части SQL с указанным префиксом.
// Если префикс уже установлен, возвращает тот же объект без копирования.
func (p *Part) WithPrefix(sql string) mrstorage.SQLPart {
	if p.sqlPrefix == sql {
		return p
	}

	c := *p
	c.sqlPrefix = sql

	return &c
}

// WithStartArg - возвращает копию части SQL с указанным номером начального аргумента.
// Если номер уже установлен, возвращает тот же объект без копирования.
func (p *Part) WithStartArg(number int) mrstorage.SQLPart {
	if p.argumentNumber == number {
		return p
	}

	c := *p
	c.argumentNumber = number

	return &c
}

// Empty - сообщает, отсутствует ли функция для формирования части SQL.
func (p *Part) Empty() bool {
	return p.partFunc == nil
}

// ToSQL - возвращает SQL-выражение в виде строки и список используемых аргументов.
func (p *Part) ToSQL() (string, []any) {
	if p.partFunc == nil {
		return "", nil
	}

	sql, args := p.partFunc(p.argumentNumber)

	return p.sqlPrefix + sql, args
}
