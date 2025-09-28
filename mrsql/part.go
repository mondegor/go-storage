package mrsql

import (
	"fmt"
	"log"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// Part - часть SQL запроса с используемыми в ней аргументами.
	Part struct {
		argumentNumber int
		sqlPrefix      string
		partFunc       mrstorage.SQLPartFunc
	}
)

// NewPart - создаёт объект Part.
func NewPart(argumentNumber int, part mrstorage.SQLPartFunc) *Part {
	return &Part{
		argumentNumber: argumentNumber,
		partFunc:       part,
	}
}

// WithPrefix - возвращает часть SQL, перед которым будет добавлен указанный префикс.
func (p *Part) WithPrefix(sql string) mrstorage.SQLPart {
	if p.sqlPrefix == sql {
		return p
	}

	c := *p
	c.sqlPrefix = sql

	return &c
}

// WithStartArg - возвращает часть SQL, в котором первый номер его аргументов будет начинаться с указанного номера.
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

// String - возвращает часть SQL в виде строки без аргументов (только если есть уверенность, что аргументы не использовались).
func (p *Part) String() string {
	sql, args := p.ToSQL()

	if len(args) > 0 {
		log.Print(fmt.Errorf("Part.String(): '%s' has %d args", sql, len(args)).Error())
	}

	return sql
}

// ToSQL - возвращает часть SQL в виде строки и отдельно используемые аргументы.
func (p *Part) ToSQL() (string, []any) {
	if p.partFunc == nil {
		return "", nil
	}

	sql, args := p.partFunc(p.argumentNumber)

	return p.sqlPrefix + sql, args
}
