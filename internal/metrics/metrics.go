package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	HttpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "endpoint"},
	)

	HttpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "endpoint"},
	)

	// Database metrics
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1.0, 5.0},
		},
		[]string{"operation", "table"},
	)

	DatabaseConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	DatabaseConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	DatabaseQueryErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_query_errors_total",
			Help: "Total number of database query errors",
		},
		[]string{"operation", "error_type"},
	)

	// Application metrics
	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of currently active users",
		},
	)

	AuthenticationAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "authentication_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"result"}, // success, failure
	)

	CacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// Business metrics
	BusinessEvents = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_events_total",
			Help: "Total number of business events",
		},
		[]string{"event_type"},
	)

	// Error metrics
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type", "severity"},
	)

	// Circuit breaker metrics
	CircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_state",
			Help: "Circuit breaker state (0=closed, 1=half-open, 2=open)",
		},
		[]string{"name"},
	)

	CircuitBreakerFailures = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_failures_total",
			Help: "Total number of circuit breaker failures",
		},
		[]string{"name"},
	)
)

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration, requestSize, responseSize int64) {
	status := strconv.Itoa(statusCode)

	HttpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	HttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	HttpRequestSize.WithLabelValues(method, endpoint).Observe(float64(requestSize))
	HttpResponseSize.WithLabelValues(method, endpoint).Observe(float64(responseSize))
}

// RecordDatabaseQuery records database query metrics
func RecordDatabaseQuery(operation, table string, duration time.Duration, err error) {
	DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())

	if err != nil {
		errorType := "unknown"
		if err != nil {
			errorType = err.Error()
		}
		DatabaseQueryErrors.WithLabelValues(operation, errorType).Inc()
	}
}

// UpdateDatabaseConnections updates database connection metrics
func UpdateDatabaseConnections(active, idle int) {
	DatabaseConnectionsActive.Set(float64(active))
	DatabaseConnectionsIdle.Set(float64(idle))
}

// RecordAuthAttempt records authentication attempt
func RecordAuthAttempt(success bool) {
	result := "failure"
	if success {
		result = "success"
	}
	AuthenticationAttempts.WithLabelValues(result).Inc()
}

// RecordCacheAccess records cache access
func RecordCacheAccess(hit bool) {
	if hit {
		CacheHits.Inc()
	} else {
		CacheMisses.Inc()
	}
}

// RecordBusinessEvent records a business event
func RecordBusinessEvent(eventType string) {
	BusinessEvents.WithLabelValues(eventType).Inc()
}

// RecordError records an error
func RecordError(errorType, severity string) {
	ErrorsTotal.WithLabelValues(errorType, severity).Inc()
}

// SetActiveUsers sets the number of active users
func SetActiveUsers(count int) {
	ActiveUsers.Set(float64(count))
}

// UpdateCircuitBreakerState updates circuit breaker metrics
func UpdateCircuitBreakerState(name string, state int) {
	CircuitBreakerState.WithLabelValues(name).Set(float64(state))
}

// RecordCircuitBreakerFailure records a circuit breaker failure
func RecordCircuitBreakerFailure(name string) {
	CircuitBreakerFailures.WithLabelValues(name).Inc()
}

// PrometheusMiddleware returns a Fiber middleware that collects metrics
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get request size
		requestSize := int64(len(c.Request().Body()))

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get response details
		statusCode := c.Response().StatusCode()
		responseSize := int64(len(c.Response().Body()))
		method := c.Method()
		path := c.Path()

		// Record metrics
		RecordHTTPRequest(method, path, statusCode, duration, requestSize, responseSize)

		return err
	}
}
