package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("cache: key not found")

type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Close() error
}

type redisCache struct {
	client *redis.Client
}

func New(redisURL string) (Cache, error) {
	if redisURL == "" {
		return &noopCache{}, nil
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &redisCache{client: client}, nil
}

func (r *redisCache) Get(ctx context.Context, key string, dest any) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return fmt.Errorf("redis get: %w", err)
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("unmarshal cached value: %w", err)
	}
	return nil
}

func (r *redisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}

func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}

func (r *redisCache) Close() error {
	return r.client.Close()
}

type noopCache struct{}

func (n *noopCache) Get(ctx context.Context, key string, dest any) error { return ErrNotFound }
func (n *noopCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return nil
}
func (n *noopCache) Del(ctx context.Context, keys ...string) error { return nil }
func (n *noopCache) Close() error                                  { return nil }
