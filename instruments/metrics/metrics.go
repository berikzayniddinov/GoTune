package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	InstrumentCreateAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "instrument_create_attempts_total",
			Help: "Total number of instrument creation attempts",
		},
	)
	InstrumentCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "instrument_created_total",
			Help: "Total number of successfully created instruments",
		},
	)
	InstrumentCreateDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "instrument_create_duration_seconds",
			Help:    "Duration of instrument creation in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	InstrumentGetAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "instrument_get_attempts_total",
			Help: "Total number of instrument get attempts",
		},
	)
	InstrumentGetDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "instrument_get_duration_seconds",
			Help:    "Duration of instrument get requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Register() {
	prometheus.MustRegister(
		InstrumentCreateAttempts,
		InstrumentCreatedTotal,
		InstrumentCreateDuration,
		InstrumentGetAttempts,
		InstrumentGetDuration,
	)
}
