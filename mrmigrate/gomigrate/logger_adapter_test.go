package gomigrate_test

import (
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrmigrate/gomigrate"
)

// Make sure the gomigrate.LoggerAdapter conforms with the migrate.Logger interface.
func TestLoggerAdapterImplementsLogger(t *testing.T) {
	assert.Implements(t, (*migrate.Logger)(nil), &gomigrate.LoggerAdapter{})
}
