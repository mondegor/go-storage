//nolint:errcheck
package placeholdedvalues

import (
	"strconv"
)

const (
	// generate: ($1), ...
	defaultCountArgs     = 1
	defaultLineStart     = "("
	defaultLineEnd       = ")"
	defaultLinePrefix    = ""
	defaultLinePostfix   = ""
	defaultArgsSeparator = ", "
	defaultLineSeparator = ", "
)

type (
	// SQL - объект позволяет формировать повторяющиеся последовательности, в которых
	// содержатся пронумерованные аргументы (например, используется, для множественной вставки в INSERT запросах).
	SQL struct {
		buf           writer
		countArgs     int
		linePrefix    string
		lineMiddle    map[int]string
		linePostfix   string
		argsSeparator string
		lineSeparator string
	}

	writer interface {
		WriteByte(value byte) error
		WriteString(value string) (int, error)
	}

	sql struct {
		p         *SQL
		lineStart string
		lineEnd   string
	}
)

// New - создаёт объект SQL.
func New(buf writer, opts ...Option) *SQL {
	ws := sql{
		p: &SQL{
			buf:           buf,
			countArgs:     defaultCountArgs,
			linePrefix:    defaultLinePrefix,
			lineMiddle:    nil,
			linePostfix:   defaultLinePostfix,
			argsSeparator: defaultArgsSeparator,
			lineSeparator: defaultLineSeparator,
		},
		lineStart: defaultLineStart,
		lineEnd:   defaultLineEnd,
	}

	for _, opt := range opts {
		opt(&ws)
	}

	// расставляются начальная и завершающая строки
	ws.p.linePrefix = ws.lineStart + ws.p.linePrefix
	ws.p.linePostfix += ws.lineEnd

	return ws.p
}

// WriteFirstLine - добавляет первую линию с аргументами.
// Параметр argumentNumber является необязательным, если он меньше или равен нулю, то он будет приравнен к 1.
// Пример: '($1, $2, $3, NOW())'.
func (s *SQL) WriteFirstLine(argumentNumber ...int) (nextArgument int) {
	if len(argumentNumber) == 0 || argumentNumber[0] < 1 {
		nextArgument = 1
	} else {
		nextArgument = argumentNumber[0]
	}

	return s.writeLine(nextArgument)
}

// WriteNextLine - добавляет запятую и следующую линию с аргументами.
// Если argumentNumber меньше или равен нулю, то он будет приравнен к 1.
// Пример: ', ($1, $2, $3, NOW())'.
func (s *SQL) WriteNextLine(argumentNumber int) (nextArgument int) {
	if argumentNumber < 1 {
		argumentNumber = 1
	}

	s.buf.WriteString(s.lineSeparator)

	return s.writeLine(argumentNumber)
}

func (s *SQL) writeLine(argumentNumber int) (nextArgumentNumber int) {
	s.buf.WriteString(s.linePrefix)

	// зная, что s.countArgs всегда > 0, последний аргумент обрабатывается отдельно
	// чтобы не использовать дополнительную проверку внутри цикла
	for i := 0; i < s.countArgs-1; i++ {
		s.buf.WriteByte('$')
		s.buf.WriteString(strconv.FormatInt(int64(argumentNumber+i), 10))

		if middle, ok := s.lineMiddle[i+1]; ok {
			s.buf.WriteString(middle)
		} else {
			s.buf.WriteString(s.argsSeparator)
		}
	}

	argumentNumber += s.countArgs - 1

	s.buf.WriteByte('$')
	s.buf.WriteString(strconv.FormatInt(int64(argumentNumber), 10))

	s.buf.WriteString(s.linePostfix)

	return argumentNumber + 1
}
