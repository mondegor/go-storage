//nolint:errcheck
package placeholdedvalues

import (
	"io"
	"strconv"
)

const (
	// generate: ($1), ...
	defaultCountArgs     = 1
	defaultLineStart     = "("
	defaultLineEnd       = ")"
	defaultArgsSeparator = ", "
	defaultLineSeparator = ", "
)

type (
	// SQL - объект позволяет формировать повторяющиеся последовательности, в которых
	// содержатся пронумерованные аргументы (например, используется, для множественной вставки в INSERT запросах).
	SQL interface {
		CountLineArgs() int
		WriteFirstLine(w io.StringWriter, argumentNumber ...int) (nextArgument int)
		WriteNextLine(w io.StringWriter, argumentNumber int) (nextArgument int)
	}

	sql struct {
		countArgs     int
		lineSpans     []string
		argsSeparator string
		lineSeparator string
	}
)

// New - создаёт объект SQL.
func New(opts ...Option) SQL {
	o := options{
		sql: &sql{
			argsSeparator: defaultArgsSeparator,
			lineSeparator: defaultLineSeparator,
		},
		lineStart: defaultLineStart,
		lineEnd:   defaultLineEnd,
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.sql.countArgs < 1 {
		o.sql.countArgs = defaultCountArgs
	}

	switch {
	// кол-во промежутков между элементами не должно быть больше o.sql.countArgs+1
	case len(o.sql.lineSpans) > o.sql.countArgs+1:
		o.sql.lineSpans = o.sql.lineSpans[: o.sql.countArgs+1 : o.sql.countArgs+1]

	// если промежутков вообще не было указано
	case len(o.sql.lineSpans) == 0:
		o.sql.lineSpans = make([]string, o.sql.countArgs+1)
	// иначе промежутки указаны, но их не больше чем аргументов
	default:
		// расширение массива промежутков до нужного кол-ва
		for i := len(o.sql.lineSpans); i < o.sql.countArgs+1; i++ {
			o.sql.lineSpans = append(o.sql.lineSpans, "")
		}
	}

	// вставка разделителя в пустые промежутки (кроме первого и последнего)
	for i := 1; i < len(o.sql.lineSpans)-1; i++ {
		if o.sql.lineSpans[i] == "" {
			o.sql.lineSpans[i] = o.sql.argsSeparator
		}
	}

	// гарантируется, что кол-во элементов в o.sql.lineSpans более 1
	o.sql.lineSpans[0] = o.lineStart + o.sql.lineSpans[0]
	o.sql.lineSpans[len(o.sql.lineSpans)-1] += o.lineEnd

	return o.sql
}

// CountLineArgs - возвращает кол-во аргументов в линии.
func (s *sql) CountLineArgs() int {
	return s.countArgs
}

// WriteFirstLine - добавляет первую линию с аргументами.
// Параметр argumentNumber является необязательным, если он меньше или равен нулю, то он будет приравнен к 1.
// Пример: '($1, $2, $3, NOW())'.
func (s *sql) WriteFirstLine(w io.StringWriter, argumentNumber ...int) (nextArgument int) {
	if len(argumentNumber) == 0 || argumentNumber[0] < 1 {
		nextArgument = 1
	} else {
		nextArgument = argumentNumber[0]
	}

	return s.writeLine(w, nextArgument)
}

// WriteNextLine - добавляет запятую и следующую линию с аргументами.
// Если argumentNumber меньше или равен нулю, то он будет приравнен к 1.
// Пример: ', ($1, $2, $3, NOW())'.
func (s *sql) WriteNextLine(w io.StringWriter, argumentNumber int) (nextArgument int) {
	if argumentNumber < 1 {
		argumentNumber = 1
	}

	w.WriteString(s.lineSeparator)

	return s.writeLine(w, argumentNumber)
}

func (s *sql) writeLine(w io.StringWriter, argumentNumber int) (nextArgumentNumber int) {
	w.WriteString(s.lineSpans[0])

	for i := 0; i < s.countArgs; i++ {
		w.WriteString("$")
		w.WriteString(strconv.FormatInt(int64(argumentNumber+i), 10))
		w.WriteString(s.lineSpans[i+1])
	}

	return argumentNumber + 1
}
