package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/itsahyarr/go-fiber-boilerplate/config"
	"github.com/redis/go-redis/v9"
)

type KeyDB struct {
	Client *redis.Client
}

func NewKeyDB(cfg config.KeyDBConfig) (*KeyDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to KeyDB: %w", err)
	}

	log.Println("✅ Connected to KeyDB successfully")

	return &KeyDB{Client: client}, nil
}

func (k *KeyDB) Close() error {
	if err := k.Client.Close(); err != nil {
		return fmt.Errorf("failed to close KeyDB connection: %w", err)
	}
	log.Println("✅ KeyDB connection closed")
	return nil
}

// Get retrieves a value from KeyDB
func (k *KeyDB) Get(ctx context.Context, key string) (string, error) {
	val, err := k.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	}
	return val, err
}

// Set stores a value in KeyDB with optional TTL
func (k *KeyDB) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return k.Client.Set(ctx, key, value, ttl).Err()
}

// Del deletes a key from KeyDB
func (k *KeyDB) Del(ctx context.Context, keys ...string) error {
	return k.Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists in KeyDB
func (k *KeyDB) Exists(ctx context.Context, keys ...string) (bool, error) {
	count, err := k.Client.Exists(ctx, keys...).Result()
	return count > 0, err
}

// TTL returns the remaining time to live of a key
func (k *KeyDB) TTL(ctx context.Context, key string) (time.Duration, error) {
	return k.Client.TTL(ctx, key).Result()
}

// GetKeys retrieves all keys matching a pattern
func (k *KeyDB) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	return k.Client.Keys(ctx, pattern).Result()
}

// SetNX sets a value only if the key does not exist (atomic operation)
func (k *KeyDB) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	return k.Client.SetNX(ctx, key, value, ttl).Result()
}

// GetDel gets and deletes a key atomically
func (k *KeyDB) GetDel(ctx context.Context, key string) (string, error) {
	return k.Client.GetDel(ctx, key).Result()
}
