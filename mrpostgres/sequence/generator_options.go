package sequence

type (
	// Option - функция для настройки объекта Generator.
	Option func(o *options)

	options struct {
		generator *Generator
	}
)

// WithMaxIDsOneQuery - устанавливает максимальное количество ID, получаемых за один запрос к БД.
// Используется методом MultiNext для оптимизации получения большого количества ID.
func WithMaxIDsOneQuery(value int) Option {
	return func(o *options) {
		o.generator.maxIDsOneQuery = value
	}
}
