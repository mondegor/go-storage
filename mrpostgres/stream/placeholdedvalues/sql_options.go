package placeholdedvalues

type (
	// Option - настройка объекта SQL.
	Option func(s *sql)
)

// WithCountArgs - устанавливает количество аргументов на одну строку (запись).
func WithCountArgs(value int) Option {
	return func(s *sql) {
		if value > 0 {
			s.p.countArgs = value
		}
	}
}

// WithLineStart - устанавливает строку, начинающую линию (до префикса).
func WithLineStart(value string) Option {
	return func(s *sql) {
		s.lineStart = value
	}
}

// WithLineEnd - устанавливает строку, завершающую линию (после постфикса).
func WithLineEnd(value string) Option {
	return func(s *sql) {
		s.lineEnd = value
	}
}

// WithLinePrefix - устанавливает префикс, который будет поставлен
// перед первым аргументом, но после начинающей скобочки.
func WithLinePrefix(value string) Option {
	return func(s *sql) {
		s.p.linePrefix = value
	}
}

// WithLineMiddle - устанавливает строки после номеров аргументов,
// где map[int]string - номер аргумента (за исключением последнего) - устанавливаемое значение сразу после этого аргумента.
// При этом нужно устанавливать запятую, разделяющие аргументы.
func WithLineMiddle(value map[int]string) Option {
	return func(s *sql) {
		s.p.lineMiddle = value
	}
}

// WithLinePostfix - устанавливает постфикс, который будет поставлен сразу
// после последнего аргумента, но до завершающей скобочки.
func WithLinePostfix(value string) Option {
	return func(s *sql) {
		s.p.linePostfix = value
	}
}

// WithArgsSeparator - устанавливает разделитель между аргументами.
// Внимание: он работает только для аргументов, которых нет в lineMiddle.
func WithArgsSeparator(value string) Option {
	return func(s *sql) {
		s.p.argsSeparator = value
	}
}

// WithLineSeparator - устанавливает разделитель между линиями.
func WithLineSeparator(value string) Option {
	return func(s *sql) {
		s.p.lineSeparator = value
	}
}
