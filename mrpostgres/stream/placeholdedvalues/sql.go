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
		countArgs     uint64
		linePrefix    string
		lineMiddle    map[uint64]string
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
// Пример: '($1, $2, $3, NOW())'.
func (s *SQL) WriteFirstLine(argumentNumber ...uint64) (nextArgument uint64) {
	if len(argumentNumber) == 0 {
		nextArgument = 1
	} else {
		nextArgument = argumentNumber[0]
	}

	return s.writeLine(nextArgument)
}

// WriteNextLine - добавляет запятую и следующую линию с аргументами.
// Пример: ', ($1, $2, $3, NOW())'.
func (s *SQL) WriteNextLine(argumentNumber uint64) (nextArgument uint64) {
	s.buf.WriteString(s.lineSeparator)

	return s.writeLine(argumentNumber)
}

func (s *SQL) writeLine(argumentNumber uint64) (nextArgumentNumber uint64) {
	s.buf.WriteString(s.linePrefix)

	// зная, что s.countArgs всегда > 0, последний аргумент обрабатывается отдельно
	// чтобы не использовать дополнительную проверку внутри цикла
	for i := uint64(0); i < s.countArgs-1; i++ {
		s.buf.WriteByte('$')
		s.buf.WriteString(strconv.FormatUint(argumentNumber+i, 10))

		if middle, ok := s.lineMiddle[i+1]; ok {
			s.buf.WriteString(middle)
		} else {
			s.buf.WriteString(s.argsSeparator)
		}
	}

	argumentNumber += s.countArgs - 1

	s.buf.WriteByte('$')
	s.buf.WriteString(strconv.FormatUint(argumentNumber, 10))

	s.buf.WriteString(s.linePostfix)

	return argumentNumber + 1
}
