package placeholdedvalues

type (
	// Option - настройка объекта SQL.
	Option func(o *options)

	options struct {
		sql       *sql
		lineStart string
		lineEnd   string
	}
)

// WithCountLineArgs - устанавливает количество аргументов на одну строку (запись).
func WithCountLineArgs(value int) Option {
	return func(o *options) {
		o.sql.countArgs = value
	}
}

// WithLineSeparator - устанавливает разделитель между сточками.
func WithLineSeparator(value string) Option {
	return func(o *options) {
		o.sql.lineSeparator = value
	}
}

// WithArgsSeparator - устанавливает разделитель между аргументами.
// Внимание: он вставляется только между аргументами, у которых промежуток пустой.
func WithArgsSeparator(value string) Option {
	return func(o *options) {
		o.sql.argsSeparator = value
	}
}

// WithLineStart - устанавливает строку, начинающую линию (до префикса).
func WithLineStart(value string) Option {
	return func(o *options) {
		o.lineStart = value
	}
}

// WithLineEnd - устанавливает строку, завершающую линию (после постфикса).
func WithLineEnd(value string) Option {
	return func(o *options) {
		o.lineEnd = value
	}
}

// WithLine - заполняет промежутки указанными значениями между аргументами
// начиная с левой стороны первого аргумента, их может быть указано на 1 больше кол-ва аргументов.
// В значениях вставляемых между элементами необходимо не забывать про разделитель lineSeparator.
// Допустимо указывать значения меньше кол-ва аргументов, тогда оставшиеся значения
// заменятся на разделитель lineSeparator где это нужно.
func WithLine(spans ...string) Option {
	return func(o *options) {
		o.sql.lineSpans = spans
	}
}
