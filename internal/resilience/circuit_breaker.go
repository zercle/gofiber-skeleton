package resilience

import (
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreaker wraps gobreaker.CircuitBreaker with configuration
type CircuitBreaker struct {
	*gobreaker.CircuitBreaker
}

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(counts gobreaker.Counts) bool
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

// NewCircuitBreaker creates a new circuit breaker with custom configuration
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        config.Name,
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
	}

	// Default ReadyToTrip function
	if config.ReadyToTrip != nil {
		settings.ReadyToTrip = config.ReadyToTrip
	} else {
		settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		}
	}

	if config.OnStateChange != nil {
		settings.OnStateChange = config.OnStateChange
	}

	return &CircuitBreaker{
		CircuitBreaker: gobreaker.NewCircuitBreaker(settings),
	}
}

// DefaultCircuitBreaker creates a circuit breaker with sensible defaults
func DefaultCircuitBreaker(name string) *CircuitBreaker {
	return NewCircuitBreaker(CircuitBreakerConfig{
		Name:        name,
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
	})
}

// DatabaseCircuitBreaker creates a circuit breaker optimized for database operations
func DatabaseCircuitBreaker() *CircuitBreaker {
	return NewCircuitBreaker(CircuitBreakerConfig{
		Name:        "database",
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     15 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip if 5 consecutive failures
			return counts.ConsecutiveFailures >= 5
		},
	})
}

// ExternalAPICircuitBreaker creates a circuit breaker for external API calls
func ExternalAPICircuitBreaker(apiName string) *CircuitBreaker {
	return NewCircuitBreaker(CircuitBreakerConfig{
		Name:        apiName,
		MaxRequests: 2,
		Interval:    60 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip if more than 50% failures in last 10 requests
			if counts.Requests < 10 {
				return false
			}
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.5
		},
	})
}

// Execute wraps a function with circuit breaker logic
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	return cb.CircuitBreaker.Execute(fn)
}

// ExecuteVoid wraps a function with no return value with circuit breaker logic
func (cb *CircuitBreaker) ExecuteVoid(fn func() error) error {
	_, err := cb.CircuitBreaker.Execute(func() (interface{}, error) {
		return nil, fn()
	})
	return err
}
