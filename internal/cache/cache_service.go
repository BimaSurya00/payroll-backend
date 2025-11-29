package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/itsahyarr/go-fiber-boilerplate/database"
)

type CacheService interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (bool, error)
	Clear(ctx context.Context, pattern string) error
}

type cacheService struct {
	keydb *database.KeyDB
}

func NewCacheService(keydb *database.KeyDB) CacheService {
	return &cacheService{keydb: keydb}
}

// Set stores a value in cache with optional TTL
func (s *cacheService) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.keydb.Set(ctx, key, data, ttl)
}

// Get retrieves and unmarshals a value from cache
func (s *cacheService) Get(ctx context.Context, key string, dest any) error {
	data, err := s.keydb.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Delete removes one or more keys from cache
func (s *cacheService) Delete(ctx context.Context, keys ...string) error {
	return s.keydb.Del(ctx, keys...)
}

// Exists checks if keys exist in cache
func (s *cacheService) Exists(ctx context.Context, keys ...string) (bool, error) {
	return s.keydb.Exists(ctx, keys...)
}

// Clear removes all keys matching a pattern (use with caution)
func (s *cacheService) Clear(ctx context.Context, pattern string) error {
	// Note: In production, you might want to implement this with SCAN instead of KEYS
	// to avoid blocking the Redis server with large keyspaces
	keys, err := s.keydb.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return s.keydb.Del(ctx, keys...)
	}
	return nil
}
