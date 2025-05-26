package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Общее количество зарегистрированных пользователей
	RegisteredUsersTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gotune_registered_users_total",
			Help: "Общее количество зарегистрированных пользователей",
		},
	)

	// Количество попыток регистрации
	RegistrationAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gotune_registration_attempts_total",
			Help: "Количество попыток регистрации",
		},
	)

	// Время обработки запроса на регистрацию
	RegistrationDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "gotune_registration_duration_seconds",
			Help:    "Длительность регистрации пользователя",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// Register регистрирует все метрики в Prometheus
func Register() {
	prometheus.MustRegister(RegisteredUsersTotal)
	prometheus.MustRegister(RegistrationAttempts)
	prometheus.MustRegister(RegistrationDuration)
}
