package mrprometheus_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrstorage/mrprometheus"
)

// Make sure the mrprometheus.DBCollector conforms with the prometheus.Collector interface.
func TestDBCollectorImplementsCollector(t *testing.T) {
	assert.Implements(t, (*prometheus.Collector)(nil), &mrprometheus.DBCollector{})
}
