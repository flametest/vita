package vredis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupTestClient(t *testing.T) (Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := NewClient(Config{Addr: mr.Addr()})
	return client, mr
}

// String operations

func TestSetAndGet(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	err := client.Set(ctx, "key1", "value1", 0)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	got, err := client.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != "value1" {
		t.Errorf("Get() = %v, want value1", got)
	}
}

func TestGetNonExistent(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	_, err := client.Get(ctx, "nonexistent")
	if err != redis.Nil {
		t.Errorf("Get() error = %v, want redis.Nil", err)
	}
}

func TestDel(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	client.Set(ctx, "key1", "value1", 0)
	client.Set(ctx, "key2", "value2", 0)

	err := client.Del(ctx, "key1", "key2")
	if err != nil {
		t.Fatalf("Del() error = %v", err)
	}

	count, _ := client.Exists(ctx, "key1", "key2")
	if count != 0 {
		t.Errorf("Exists() = %v, want 0", count)
	}
}

func TestExists(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	count, err := client.Exists(ctx, "key1")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if count != 0 {
		t.Errorf("Exists() = %v, want 0", count)
	}

	client.Set(ctx, "key1", "value1", 0)
	count, _ = client.Exists(ctx, "key1")
	if count != 1 {
		t.Errorf("Exists() = %v, want 1", count)
	}
}

func TestExpireAndTTL(t *testing.T) {
	client, mr := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	client.Set(ctx, "key1", "value1", 0)

	err := client.Expire(ctx, "key1", 10*time.Second)
	if err != nil {
		t.Fatalf("Expire() error = %v", err)
	}

	ttl, err := client.TTL(ctx, "key1")
	if err != nil {
		t.Fatalf("TTL() error = %v", err)
	}
	if ttl <= 0 || ttl > 10*time.Second {
		t.Errorf("TTL() = %v, want between 0 and 10s", ttl)
	}

	mr.FastForward(5 * time.Second)
	ttl, _ = client.TTL(ctx, "key1")
	if ttl > 5*time.Second {
		t.Errorf("TTL() after fast forward = %v, want <= 5s", ttl)
	}
}

func TestIncrDecr(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	val, err := client.Incr(ctx, "counter")
	if err != nil {
		t.Fatalf("Incr() error = %v", err)
	}
	if val != 1 {
		t.Errorf("Incr() = %v, want 1", val)
	}

	val, _ = client.Incr(ctx, "counter")
	if val != 2 {
		t.Errorf("Incr() = %v, want 2", val)
	}

	val, err = client.Decr(ctx, "counter")
	if err != nil {
		t.Fatalf("Decr() error = %v", err)
	}
	if val != 1 {
		t.Errorf("Decr() = %v, want 1", val)
	}
}

func TestIncrBy(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	val, err := client.IncrBy(ctx, "counter", 10)
	if err != nil {
		t.Fatalf("IncrBy() error = %v", err)
	}
	if val != 10 {
		t.Errorf("IncrBy() = %v, want 10", val)
	}

	val, _ = client.IncrBy(ctx, "counter", -3)
	if val != 7 {
		t.Errorf("IncrBy() = %v, want 7", val)
	}
}

// Hash operations

func TestHashSetAndGet(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	err := client.HSet(ctx, "hash1", "field1", "value1", "field2", "value2")
	if err != nil {
		t.Fatalf("HSet() error = %v", err)
	}

	got, err := client.HGet(ctx, "hash1", "field1")
	if err != nil {
		t.Fatalf("HGet() error = %v", err)
	}
	if got != "value1" {
		t.Errorf("HGet() = %v, want value1", got)
	}

	all, err := client.HGetAll(ctx, "hash1")
	if err != nil {
		t.Fatalf("HGetAll() error = %v", err)
	}
	if len(all) != 2 || all["field1"] != "value1" || all["field2"] != "value2" {
		t.Errorf("HGetAll() = %v, want map[field1:value1 field2:value2]", all)
	}
}

func TestHDel(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	client.HSet(ctx, "hash1", "field1", "value1")
	client.HDel(ctx, "hash1", "field1")

	exists, _ := client.HExists(ctx, "hash1", "field1")
	if exists {
		t.Error("HExists() = true, want false after HDel")
	}
}

func TestHExists(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	client.HSet(ctx, "hash1", "field1", "value1")

	exists, err := client.HExists(ctx, "hash1", "field1")
	if err != nil {
		t.Fatalf("HExists() error = %v", err)
	}
	if !exists {
		t.Error("HExists() = false, want true")
	}

	exists, _ = client.HExists(ctx, "hash1", "field2")
	if exists {
		t.Error("HExists() = true, want false")
	}
}

func TestHIncrBy(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	val, err := client.HIncrBy(ctx, "hash1", "count", 5)
	if err != nil {
		t.Fatalf("HIncrBy() error = %v", err)
	}
	if val != 5 {
		t.Errorf("HIncrBy() = %v, want 5", val)
	}

	val, _ = client.HIncrBy(ctx, "hash1", "count", 3)
	if val != 8 {
		t.Errorf("HIncrBy() = %v, want 8", val)
	}
}

// Set operations

func TestSetOps(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	err := client.SAdd(ctx, "set1", "a", "b", "c")
	if err != nil {
		t.Fatalf("SAdd() error = %v", err)
	}

	members, err := client.SMembers(ctx, "set1")
	if err != nil {
		t.Fatalf("SMembers() error = %v", err)
	}
	if len(members) != 3 {
		t.Errorf("SMembers() = %v, want 3 members", members)
	}

	isMember, _ := client.SIsMember(ctx, "set1", "a")
	if !isMember {
		t.Error("SIsMember() = false, want true")
	}

	client.SRem(ctx, "set1", "a")
	isMember, _ = client.SIsMember(ctx, "set1", "a")
	if isMember {
		t.Error("SIsMember() after SRem = true, want false")
	}
}

// Sorted Set operations

func TestSortedSetOps(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	err := client.ZAdd(ctx, "zset1",
		redis.Z{Score: 1, Member: "a"},
		redis.Z{Score: 2, Member: "b"},
		redis.Z{Score: 3, Member: "c"},
	)
	if err != nil {
		t.Fatalf("ZAdd() error = %v", err)
	}

	card, err := client.ZCard(ctx, "zset1")
	if err != nil {
		t.Fatalf("ZCard() error = %v", err)
	}
	if card != 3 {
		t.Errorf("ZCard() = %v, want 3", card)
	}

	members, err := client.ZRange(ctx, "zset1", 0, -1)
	if err != nil {
		t.Fatalf("ZRange() error = %v", err)
	}
	if len(members) != 3 || members[0] != "a" || members[2] != "c" {
		t.Errorf("ZRange() = %v, want [a b c]", members)
	}

	score, err := client.ZScore(ctx, "zset1", "b")
	if err != nil {
		t.Fatalf("ZScore() error = %v", err)
	}
	if score != 2 {
		t.Errorf("ZScore() = %v, want 2", score)
	}

	rank, err := client.ZRank(ctx, "zset1", "c")
	if err != nil {
		t.Fatalf("ZRank() error = %v", err)
	}
	if rank != 2 {
		t.Errorf("ZRank() = %v, want 2", rank)
	}

	client.ZRem(ctx, "zset1", "b")
	card, _ = client.ZCard(ctx, "zset1")
	if card != 2 {
		t.Errorf("ZCard() after ZRem = %v, want 2", card)
	}
}

func TestZRangeByScore(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	client.ZAdd(ctx, "zset1",
		redis.Z{Score: 1, Member: "a"},
		redis.Z{Score: 5, Member: "b"},
		redis.Z{Score: 10, Member: "c"},
	)

	members, err := client.ZRangeByScore(ctx, "zset1", &redis.ZRangeBy{
		Min: "2",
		Max: "8",
	})
	if err != nil {
		t.Fatalf("ZRangeByScore() error = %v", err)
	}
	if len(members) != 1 || members[0] != "b" {
		t.Errorf("ZRangeByScore() = %v, want [b]", members)
	}
}

// Pipeline

func TestPipeline(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()
	ctx := context.Background()

	err := client.Pipeline(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, "k1", "v1", 0)
		pipe.Set(ctx, "k2", "v2", 0)
		return nil
	})
	if err != nil {
		t.Fatalf("Pipeline() error = %v", err)
	}

	v1, _ := client.Get(ctx, "k1")
	v2, _ := client.Get(ctx, "k2")
	if v1 != "v1" || v2 != "v2" {
		t.Errorf("Pipeline results: k1=%v, k2=%v, want v1, v2", v1, v2)
	}
}

// Redis() accessor

func TestRedisAccessor(t *testing.T) {
	client, _ := setupTestClient(t)
	defer client.Close()

	rdb := client.Redis()
	if rdb == nil {
		t.Error("Redis() returned nil")
	}
}
