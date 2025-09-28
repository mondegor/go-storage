package sequence

type (
	// Option - настройка объекта Generator.
	Option func(g *Generator)
)

// WithMaxIDsOneQuery - устанавливает максимально возможное
// получение ID из последовательности за один запрос к БД.
func WithMaxIDsOneQuery(value int) Option {
	return func(g *Generator) {
		g.maxIDsOneQuery = value
	}
}
