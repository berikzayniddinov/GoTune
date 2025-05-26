package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CartCreateAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cart_create_attempts_total",
			Help: "Total number of attempts to create a cart",
		},
	)
	CartCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cart_created_total",
			Help: "Total number of carts successfully created",
		},
	)
	CartCreateDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cart_create_duration_seconds",
			Help:    "Duration in seconds of cart creation",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
	)

	CartGetAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cart_get_attempts_total",
			Help: "Total number of attempts to get a cart",
		},
	)
	CartGetDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cart_get_duration_seconds",
			Help:    "Duration in seconds of getting a cart",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
	)

	CartDeleteAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cart_delete_attempts_total",
			Help: "Total number of attempts to delete a cart",
		},
	)
	CartDeleteDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cart_delete_duration_seconds",
			Help:    "Duration in seconds of deleting a cart",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
	)
)

func Register() {
	prometheus.MustRegister(
		CartCreateAttempts,
		CartCreatedTotal,
		CartCreateDuration,
		CartGetAttempts,
		CartGetDuration,
		CartDeleteAttempts,
		CartDeleteDuration,
	)
}
