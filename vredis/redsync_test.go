package vredis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestLockAndUnlock(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()

	lock := NewLock(client, "test:lock:1", LockOption{
		Expiry:     5 * time.Second,
		Tries:      3,
		RetryDelay: 100 * time.Millisecond,
	})

	err := lock.Lock()
	if err != nil {
		t.Fatalf("Lock() error = %v", err)
	}

	ok, err := lock.Unlock()
	if err != nil {
		t.Fatalf("Unlock() error = %v", err)
	}
	if !ok {
		t.Error("Unlock() = false, want true")
	}
}

func TestLockConflict(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()

	lock1 := NewLock(client, "test:lock:conflict", LockOption{
		Expiry:     5 * time.Second,
		Tries:      1,
		RetryDelay: 100 * time.Millisecond,
	})
	lock2 := NewLock(client, "test:lock:conflict", LockOption{
		Expiry:     5 * time.Second,
		Tries:      1,
		RetryDelay: 100 * time.Millisecond,
	})

	if err := lock1.Lock(); err != nil {
		t.Fatalf("lock1.Lock() error = %v", err)
	}

	err := lock2.Lock()
	if err == nil {
		t.Error("lock2.Lock() should fail when lock1 holds the lock")
		lock2.Unlock()
	}

	lock1.Unlock()
}

func TestLockExtend(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()

	lock := NewLock(client, "test:lock:extend", LockOption{
		Expiry:     3 * time.Second,
		Tries:      3,
		RetryDelay: 100 * time.Millisecond,
	})

	if err := lock.Lock(); err != nil {
		t.Fatalf("Lock() error = %v", err)
	}
	defer lock.Unlock()

	ok, err := lock.Extend()
	if err != nil {
		t.Fatalf("Extend() error = %v", err)
	}
	if !ok {
		t.Error("Extend() = false, want true")
	}
}

func TestLockFailFast(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()

	lock1 := NewLock(client, "test:lock:failfast", LockOption{
		Expiry:     5 * time.Second,
		Tries:      10,
		RetryDelay: 100 * time.Millisecond,
	})
	lock2 := NewLock(client, "test:lock:failfast", LockOption{
		Expiry:     5 * time.Second,
		Tries:      10,
		RetryDelay: 100 * time.Millisecond,
		FailFast:   true,
	})

	lock1.Lock()

	start := time.Now()
	err := lock2.Lock()
	elapsed := time.Since(start)

	if err == nil {
		t.Error("lock2.Lock() should fail")
		lock2.Unlock()
	}
	if elapsed > 2*time.Second {
		t.Errorf("FailFast took %v, should have returned fast", elapsed)
	}

	lock1.Unlock()
}

func TestLockExpiredAndReacquire(t *testing.T) {
	client, mr := setupTestClientWithMiniredis(t)
	defer client.Close()

	lock := NewLock(client, "test:lock:expire", LockOption{
		Expiry:     1 * time.Second,
		Tries:      3,
		RetryDelay: 100 * time.Millisecond,
	})

	if err := lock.Lock(); err != nil {
		t.Fatalf("Lock() error = %v", err)
	}

	mr.FastForward(2 * time.Second)

	lock2 := NewLock(client, "test:lock:expire", LockOption{
		Expiry:     5 * time.Second,
		Tries:      3,
		RetryDelay: 100 * time.Millisecond,
	})
	if err := lock2.Lock(); err != nil {
		t.Fatalf("Lock() after expiry error = %v", err)
	}
	lock2.Unlock()
}

func TestUnlockNotHeld(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()

	lock := NewLock(client, "test:lock:unheld", LockOption{
		Expiry:     5 * time.Second,
		Tries:      1,
		RetryDelay: 100 * time.Millisecond,
	})

	ok, err := lock.Unlock()
	if ok {
		t.Error("Unlock() = true for non-held lock, want false")
	}
	if err == nil {
		t.Error("Unlock() should return error for non-held lock")
	}
}

func setupTestClientWithMiniredis(t *testing.T) (Client, *miniredis.Miniredis) {
	t.Helper()
	return setupTestClient(t)
}
