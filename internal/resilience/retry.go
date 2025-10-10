package resilience

import (
	"context"
	"time"

	"github.com/avast/retry-go/v4"
)

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	Attempts      uint
	Delay         time.Duration
	MaxDelay      time.Duration
	DelayType     retry.DelayTypeFunc
	OnRetry       func(n uint, err error)
	RetryIf       func(err error) bool
	Context       context.Context
	LastErrorOnly bool
}

// RetryWithBackoff retries an operation with exponential backoff
func RetryWithBackoff(operation func() error, config RetryConfig) error {
	var opts []retry.Option

	// Set attempts (default: 3)
	if config.Attempts == 0 {
		config.Attempts = 3
	}
	opts = append(opts, retry.Attempts(config.Attempts))

	// Set initial delay (default: 1 second)
	if config.Delay == 0 {
		config.Delay = 1 * time.Second
	}
	opts = append(opts, retry.Delay(config.Delay))

	// Set max delay (default: 30 seconds)
	if config.MaxDelay == 0 {
		config.MaxDelay = 30 * time.Second
	}
	opts = append(opts, retry.MaxDelay(config.MaxDelay))

	// Set delay type (default: exponential backoff)
	if config.DelayType != nil {
		opts = append(opts, retry.DelayType(config.DelayType))
	} else {
		opts = append(opts, retry.DelayType(retry.BackOffDelay))
	}

	// Set retry condition
	if config.RetryIf != nil {
		opts = append(opts, retry.RetryIf(config.RetryIf))
	}

	// Set context
	if config.Context != nil {
		opts = append(opts, retry.Context(config.Context))
	}

	// Set callback
	if config.OnRetry != nil {
		opts = append(opts, retry.OnRetry(config.OnRetry))
	}

	// Return only last error
	if config.LastErrorOnly {
		opts = append(opts, retry.LastErrorOnly(true))
	}

	return retry.Do(operation, opts...)
}

// SimpleRetry retries an operation with fixed delay and default settings
func SimpleRetry(operation func() error, attempts uint) error {
	return RetryWithBackoff(operation, RetryConfig{
		Attempts:      attempts,
		Delay:         1 * time.Second,
		MaxDelay:      5 * time.Second,
		DelayType:     retry.FixedDelay,
		LastErrorOnly: true,
	})
}

// DatabaseRetry retries database operations with appropriate settings
func DatabaseRetry(operation func() error) error {
	return RetryWithBackoff(operation, RetryConfig{
		Attempts:  3,
		Delay:     500 * time.Millisecond,
		MaxDelay:  5 * time.Second,
		DelayType: retry.BackOffDelay,
		OnRetry: func(n uint, err error) {
			// Log retry attempt (in production, use structured logging)
			// log.Printf("Database operation retry %d: %v", n, err)
		},
		RetryIf: func(err error) bool {
			// Only retry on transient errors
			// Check for connection errors, deadlocks, etc.
			return isTransientDatabaseError(err)
		},
		LastErrorOnly: true,
	})
}

// ExternalAPIRetry retries external API calls with appropriate settings
func ExternalAPIRetry(operation func() error) error {
	return RetryWithBackoff(operation, RetryConfig{
		Attempts:  5,
		Delay:     2 * time.Second,
		MaxDelay:  30 * time.Second,
		DelayType: retry.BackOffDelay,
		OnRetry: func(n uint, err error) {
			// Log retry attempt
			// log.Printf("External API retry %d: %v", n, err)
		},
		LastErrorOnly: true,
	})
}

// ContextAwareRetry retries an operation with context cancellation support
func ContextAwareRetry(ctx context.Context, operation func() error, attempts uint) error {
	return RetryWithBackoff(operation, RetryConfig{
		Attempts:      attempts,
		Delay:         1 * time.Second,
		MaxDelay:      10 * time.Second,
		DelayType:     retry.BackOffDelay,
		Context:       ctx,
		LastErrorOnly: true,
	})
}

// isTransientDatabaseError checks if an error is transient and worth retrying
func isTransientDatabaseError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific transient error patterns
	// This is a simplified check; in production, inspect actual error types
	errStr := err.Error()

	// Common transient errors
	transientPatterns := []string{
		"connection refused",
		"connection reset",
		"broken pipe",
		"deadlock",
		"timeout",
		"temporary failure",
	}

	for _, pattern := range transientPatterns {
		if contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(str, substr string) bool {
	// Simple case-insensitive check
	return len(str) >= len(substr) &&
		   (str == substr ||
		    stringContains(str, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
