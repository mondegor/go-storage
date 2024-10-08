package mrprometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/mondegor/go-storage/mrstorage"
)

// go get -u github.com/prometheus/client_golang

// DBCollector - представляет собой prometheus.Collector,
// который собирает статистику получаемую из DBStatProvider.
type (
	DBCollector struct {
		providerFunc func() mrstorage.DBStatProvider

		acquireCountDesc            *prometheus.Desc
		acquireDurationDesc         *prometheus.Desc
		acquiredConnsDesc           *prometheus.Desc
		canceledAcquireCountDesc    *prometheus.Desc
		constructingConnsDesc       *prometheus.Desc
		emptyAcquireCountDesc       *prometheus.Desc
		idleConnsDesc               *prometheus.Desc
		maxConnsDesc                *prometheus.Desc
		totalConnsDesc              *prometheus.Desc
		newConnsCountDesc           *prometheus.Desc
		maxLifetimeDestroyCountDesc *prometheus.Desc
		maxIdleDestroyCountDesc     *prometheus.Desc
	}
)

// NewDBCollector - создаёт объект DBCollector.
func NewDBCollector(namespace string, providerFunc func() mrstorage.DBStatProvider, labels map[string]string) *DBCollector {
	fqName := func(name string) string {
		return prometheus.BuildFQName(namespace, "pool", name)
	}

	return &DBCollector{
		providerFunc: providerFunc,

		acquireCountDesc: prometheus.NewDesc(
			fqName("acquire_count"),
			"Cumulative count of successful acquires from the pool.",
			nil,
			labels,
		),
		acquireDurationDesc: prometheus.NewDesc(
			fqName("acquire_duration_ns"),
			"Total duration of all successful acquires from the pool in nanoseconds.",
			nil,
			labels,
		),
		acquiredConnsDesc: prometheus.NewDesc(
			fqName("acquired_conns"),
			"Number of currently acquired connections in the pool.",
			nil,
			labels,
		),
		canceledAcquireCountDesc: prometheus.NewDesc(
			fqName("canceled_acquire_count"),
			"Cumulative count of acquires from the pool that were canceled by a context.",
			nil,
			labels,
		),
		constructingConnsDesc: prometheus.NewDesc(
			fqName("constructing_conns"),
			"Number of conns with construction in progress in the pool.",
			nil,
			labels,
		),
		emptyAcquireCountDesc: prometheus.NewDesc(
			fqName("empty_acquire"),
			"Cumulative count of successful acquires from the pool that waited for a resource to be released or constructed because the pool was empty.",
			nil,
			labels,
		),
		idleConnsDesc: prometheus.NewDesc(
			fqName("idle_conns"),
			"Number of currently idle conns in the pool.",
			nil,
			labels,
		),
		maxConnsDesc: prometheus.NewDesc(
			fqName("max_conns"),
			"Maximum size of the pool.",
			nil,
			labels,
		),
		totalConnsDesc: prometheus.NewDesc(
			fqName("total_conns"),
			"Total number of resources currently in the pool. The value is the sum of ConstructingConns, AcquiredConns, and IdleConns.",
			nil,
			labels,
		),
		newConnsCountDesc: prometheus.NewDesc(
			fqName("new_conns_count"),
			"Cumulative count of new connections opened.",
			nil, labels,
		),
		maxLifetimeDestroyCountDesc: prometheus.NewDesc(
			fqName("max_lifetime_destroy_count"),
			"Cumulative count of connections destroyed because they exceeded MaxConnLifetime.",
			nil,
			labels,
		),
		maxIdleDestroyCountDesc: prometheus.NewDesc(
			fqName("max_idle_destroy_count"),
			"Cumulative count of connections destroyed because they exceeded MaxConnIdleTime.",
			nil,
			labels,
		),
	}
}

// Describe - implements the prometheus.Collector interface.
func (c *DBCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect implements the prometheus.Collector interface.
func (c *DBCollector) Collect(metrics chan<- prometheus.Metric) {
	statProvider := c.providerFunc()

	sendConstMetric := func(desc *prometheus.Desc, valueType prometheus.ValueType, value float64) {
		metrics <- prometheus.MustNewConstMetric(desc, valueType, value)
	}

	sendConstMetric(c.acquireCountDesc, prometheus.CounterValue, float64(statProvider.AcquireCount()))
	sendConstMetric(c.acquireDurationDesc, prometheus.CounterValue, float64(statProvider.AcquireDuration()))
	sendConstMetric(c.acquiredConnsDesc, prometheus.GaugeValue, float64(statProvider.AcquiredConns()))
	sendConstMetric(c.canceledAcquireCountDesc, prometheus.CounterValue, float64(statProvider.CanceledAcquireCount()))
	sendConstMetric(c.constructingConnsDesc, prometheus.GaugeValue, float64(statProvider.ConstructingConns()))
	sendConstMetric(c.emptyAcquireCountDesc, prometheus.CounterValue, float64(statProvider.EmptyAcquireCount()))
	sendConstMetric(c.idleConnsDesc, prometheus.GaugeValue, float64(statProvider.IdleConns()))
	sendConstMetric(c.maxConnsDesc, prometheus.GaugeValue, float64(statProvider.MaxConns()))
	sendConstMetric(c.totalConnsDesc, prometheus.GaugeValue, float64(statProvider.TotalConns()))
	sendConstMetric(c.newConnsCountDesc, prometheus.CounterValue, float64(statProvider.NewConnsCount()))
	sendConstMetric(c.maxLifetimeDestroyCountDesc, prometheus.CounterValue, float64(statProvider.MaxLifetimeDestroyCount()))
	sendConstMetric(c.maxIdleDestroyCountDesc, prometheus.CounterValue, float64(statProvider.MaxIdleDestroyCount()))
}
