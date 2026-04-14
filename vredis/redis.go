package vredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr           string `json:"addr" yaml:"addr"`
	Password       string `json:"password" yaml:"password"`
	DB             int    `json:"db" yaml:"db"`
	PoolSize       int    `json:"poolSize" yaml:"poolSize"`
	MaxActiveConns int    `json:"maxActiveConns" yaml:"maxActiveConns"`
	MaxIdleConns   int    `json:"maxIdleConns" yaml:"maxIdleConns"`
	MinIdleConns   int    `json:"minIdleConns" yaml:"minIdleConns"`
	DialTimeout    int    `json:"dialTimeout" yaml:"dialTimeout"`
	ReadTimeout    int    `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout   int    `json:"writeTimeout" yaml:"writeTimeout"`
}

type Client interface {
	// String
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)

	// Hash
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HDel(ctx context.Context, key string, fields ...string) error
	HExists(ctx context.Context, key, field string) (bool, error)
	HIncrBy(ctx context.Context, key, field string, value int64) (int64, error)

	// Set
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)

	// Sorted Set
	ZAdd(ctx context.Context, key string, members ...redis.Z) error
	ZRem(ctx context.Context, key string, members ...interface{}) error
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZScore(ctx context.Context, key, member string) (float64, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZRank(ctx context.Context, key, member string) (int64, error)

	// Advanced
	Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) error
	Redis() *redis.Client
	Close() error
}

type redisClientImpl struct {
	client *redis.Client
}

func NewClient(config Config) Client {
	opts := &redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	}
	if config.MaxActiveConns > 0 {
		opts.MaxActiveConns = config.MaxActiveConns
	}
	if config.MaxIdleConns > 0 {
		opts.MaxIdleConns = config.MaxIdleConns
	}
	if config.DialTimeout > 0 {
		opts.DialTimeout = time.Duration(config.DialTimeout) * time.Second
	}
	if config.ReadTimeout > 0 {
		opts.ReadTimeout = time.Duration(config.ReadTimeout) * time.Second
	}
	if config.WriteTimeout > 0 {
		opts.WriteTimeout = time.Duration(config.WriteTimeout) * time.Second
	}

	return &redisClientImpl{client: redis.NewClient(opts)}
}

// String operations

func (r *redisClientImpl) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisClientImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisClientImpl) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisClientImpl) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

func (r *redisClientImpl) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r *redisClientImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *redisClientImpl) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *redisClientImpl) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

func (r *redisClientImpl) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// Hash operations

func (r *redisClientImpl) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

func (r *redisClientImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

func (r *redisClientImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

func (r *redisClientImpl) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

func (r *redisClientImpl) HExists(ctx context.Context, key, field string) (bool, error) {
	return r.client.HExists(ctx, key, field).Result()
}

func (r *redisClientImpl) HIncrBy(ctx context.Context, key, field string, value int64) (int64, error) {
	return r.client.HIncrBy(ctx, key, field, value).Result()
}

// Set operations

func (r *redisClientImpl) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

func (r *redisClientImpl) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

func (r *redisClientImpl) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

func (r *redisClientImpl) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

// Sorted Set operations

func (r *redisClientImpl) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return r.client.ZAdd(ctx, key, members...).Err()
}

func (r *redisClientImpl) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.ZRem(ctx, key, members...).Err()
}

func (r *redisClientImpl) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}

func (r *redisClientImpl) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRangeByScore(ctx, key, opt).Result()
}

func (r *redisClientImpl) ZScore(ctx context.Context, key, member string) (float64, error) {
	return r.client.ZScore(ctx, key, member).Result()
}

func (r *redisClientImpl) ZCard(ctx context.Context, key string) (int64, error) {
	return r.client.ZCard(ctx, key).Result()
}

func (r *redisClientImpl) ZRank(ctx context.Context, key, member string) (int64, error) {
	return r.client.ZRank(ctx, key, member).Result()
}

// Advanced

func (r *redisClientImpl) Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) error {
	pipe := r.client.Pipeline()
	if err := fn(pipe); err != nil {
		return err
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisClientImpl) Redis() *redis.Client {
	return r.client
}

func (r *redisClientImpl) Close() error {
	return r.client.Close()
}
