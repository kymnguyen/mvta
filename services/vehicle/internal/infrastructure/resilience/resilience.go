package resilience

import (
	"context"
	"fmt"
	"time"
)

// CircuitBreakerState represents the state of a circuit breaker.
type CircuitBreakerState string

const (
	StateClosed   CircuitBreakerState = "closed"
	StateOpen     CircuitBreakerState = "open"
	StateHalfOpen CircuitBreakerState = "half-open"
)

// CircuitBreaker implements the circuit breaker pattern for resilience.
type CircuitBreaker struct {
	state            CircuitBreakerState
	failureCount     int
	successCount     int
	failureThreshold int
	successThreshold int
	timeout          time.Duration
	lastFailureTime  time.Time
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureCount:     0,
		successCount:     0,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	}
}

// Execute runs a function with circuit breaker protection.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	switch cb.state {
	case StateClosed:
		return cb.executeClosed(fn)
	case StateOpen:
		return cb.executeOpen()
	case StateHalfOpen:
		return cb.executeHalfOpen(fn)
	default:
		return fmt.Errorf("unknown circuit breaker state: %s", cb.state)
	}
}

func (cb *CircuitBreaker) executeClosed(fn func() error) error {
	if err := fn(); err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		if cb.failureCount >= cb.failureThreshold {
			cb.state = StateOpen
			cb.failureCount = 0
		}
		return err
	}
	cb.failureCount = 0
	return nil
}

func (cb *CircuitBreaker) executeOpen() error {
	if time.Since(cb.lastFailureTime) > cb.timeout {
		cb.state = StateHalfOpen
		cb.successCount = 0
		return nil
	}
	return fmt.Errorf("circuit breaker is open, request rejected")
}

func (cb *CircuitBreaker) executeHalfOpen(fn func() error) error {
	if err := fn(); err != nil {
		cb.state = StateOpen
		cb.lastFailureTime = time.Now()
		cb.successCount = 0
		return err
	}

	cb.successCount++
	if cb.successCount >= cb.successThreshold {
		cb.state = StateClosed
		cb.failureCount = 0
		cb.successCount = 0
	}
	return nil
}

// RetryPolicy defines retry behavior.
type RetryPolicy struct {
	maxAttempts int
	backoff     time.Duration
	maxBackoff  time.Duration
}

// NewRetryPolicy creates a new retry policy.
func NewRetryPolicy(maxAttempts int, backoff, maxBackoff time.Duration) *RetryPolicy {
	return &RetryPolicy{
		maxAttempts: maxAttempts,
		backoff:     backoff,
		maxBackoff:  maxBackoff,
	}
}

// Execute runs a function with retry logic.
func (rp *RetryPolicy) Execute(ctx context.Context, fn func() error) error {
	var lastErr error
	backoff := rp.backoff

	for attempt := 0; attempt < rp.maxAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return ctx.Err()
			}
			// Exponential backoff with cap
			backoff *= 2
			if backoff > rp.maxBackoff {
				backoff = rp.maxBackoff
			}
		}

		if err := fn(); err != nil {
			lastErr = err
			continue
		}
		return nil
	}

	return fmt.Errorf("operation failed after %d attempts: %w", rp.maxAttempts, lastErr)
}

// IdempotencyKey ensures operations are idempotent.
type IdempotencyKey struct {
	key string
}

// NewIdempotencyKey creates a new idempotency key.
func NewIdempotencyKey(key string) IdempotencyKey {
	return IdempotencyKey{key: key}
}

// Key returns the idempotency key string.
func (ik IdempotencyKey) Key() string {
	return ik.key
}

// IdempotencyStore manages idempotent operation results.
type IdempotencyStore interface {
	// Store stores the result of an idempotent operation.
	Store(ctx context.Context, key IdempotencyKey, result interface{}) error

	// Get retrieves a stored result for an idempotent operation.
	Get(ctx context.Context, key IdempotencyKey) (interface{}, bool, error)
}
