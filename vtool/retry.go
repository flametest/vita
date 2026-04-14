package vtool

import (
	"context"
	"math"
	"time"
)

type RetryOption struct {
	// MaxAttempts is the maximum number of attempts (including the first call).
	MaxAttempts int
	// InitialInterval is the wait time before the first retry.
	InitialInterval time.Duration
	// MaxInterval caps the wait time to prevent unbounded growth.
	MaxInterval time.Duration
	// Multiplier controls the exponential growth rate. Defaults to 2.
	Multiplier float64
}

// Retry executes fn repeatedly until it succeeds, the context is cancelled,
// or MaxAttempts is reached. The interval between retries grows exponentially.
//
// The wait time formula is: min(InitialInterval * Multiplier^attempt, MaxInterval)
func Retry(ctx context.Context, fn func() error, opt RetryOption) error {
	if opt.MaxAttempts <= 0 {
		opt.MaxAttempts = 1
	}
	if opt.InitialInterval <= 0 {
		opt.InitialInterval = time.Second
	}
	if opt.MaxInterval <= 0 {
		opt.MaxInterval = 30 * time.Second
	}
	if opt.Multiplier <= 0 {
		opt.Multiplier = 2
	}

	var err error
	for attempt := 0; attempt < opt.MaxAttempts; attempt++ {
		if err = fn(); err == nil {
			return nil
		}
		if attempt == opt.MaxAttempts-1 {
			break
		}

		wait := time.Duration(float64(opt.InitialInterval) * math.Pow(opt.Multiplier, float64(attempt)))
		if wait > opt.MaxInterval {
			wait = opt.MaxInterval
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
	return err
}
