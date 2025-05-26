package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	OrderCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "order_created_total",
			Help: "Общее количество успешно созданных заказов",
		},
	)

	OrderCreateAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "order_create_attempts_total",
			Help: "Общее количество попыток создания заказа",
		},
	)

	OrderCreateDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "order_create_duration_seconds",
			Help:    "Продолжительность обработки создания заказа",
			Buckets: prometheus.DefBuckets,
		},
	)

	OrderGetAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "order_get_attempts_total",
			Help: "Количество запросов на получение заказа",
		},
	)

	OrderGetDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "order_get_duration_seconds",
			Help:    "Продолжительность получения информации о заказе",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Register() {
	prometheus.MustRegister(
		OrderCreatedTotal,
		OrderCreateAttempts,
		OrderCreateDuration,
		OrderGetAttempts,
		OrderGetDuration,
	)
}
