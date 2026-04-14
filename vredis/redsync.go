package vredis

import (
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

type LockOption struct {
	// Expiry is the lock's auto-expire time. Default is 8s.
	Expiry time.Duration
	// Tries is the number of attempts to acquire the lock. Default is 32.
	Tries int
	// RetryDelay is the delay between retry attempts. Default is 500ms.
	RetryDelay time.Duration
	// DriftFactor is the tolerance for clock drift between nodes. Default is 0.01.
	DriftFactor float64
	// TimeoutFactor is the timeout factor for each attempt. Default is 0.5.
	TimeoutFactor float64
	// FailFast skips retry and returns immediately on first failure. Default is false.
	FailFast bool
}

type Lock interface {
	// Lock acquires the distributed lock. Blocks until acquired or max retries reached.
	Lock() error
	// Unlock releases the distributed lock.
	Unlock() (bool, error)
	// Extend renews the lock's expiry time.
	Extend() (bool, error)
}

type distributedLock struct {
	mutex *redsync.Mutex
}

// NewLock creates a distributed lock for the given key using the provided Redis client.
func NewLock(client Client, key string, opt LockOption) Lock {
	pool := goredis.NewPool(client.Redis())
	rs := redsync.New(pool)

	var opts []redsync.Option
	if opt.Expiry > 0 {
		opts = append(opts, redsync.WithExpiry(opt.Expiry))
	}
	if opt.Tries > 0 {
		opts = append(opts, redsync.WithTries(opt.Tries))
	}
	if opt.RetryDelay > 0 {
		opts = append(opts, redsync.WithRetryDelay(opt.RetryDelay))
	}
	if opt.DriftFactor > 0 {
		opts = append(opts, redsync.WithDriftFactor(opt.DriftFactor))
	}
	if opt.TimeoutFactor > 0 {
		opts = append(opts, redsync.WithTimeoutFactor(opt.TimeoutFactor))
	}
	if opt.FailFast {
		opts = append(opts, redsync.WithFailFast(true))
	}

	return &distributedLock{
		mutex: rs.NewMutex(key, opts...),
	}
}

func (l *distributedLock) Lock() error {
	return l.mutex.Lock()
}

func (l *distributedLock) Unlock() (bool, error) {
	return l.mutex.Unlock()
}

func (l *distributedLock) Extend() (bool, error) {
	return l.mutex.Extend()
}
