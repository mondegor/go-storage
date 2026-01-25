package sequence

type (
	// Option - настройка объекта Generator.
	Option func(o *options)

	options struct {
		generator *Generator
	}
)

// WithMaxIDsOneQuery - устанавливает максимально возможное
// получение ID из последовательности за один запрос к БД.
func WithMaxIDsOneQuery(value int) Option {
	return func(o *options) {
		o.generator.maxIDsOneQuery = value
	}
}
