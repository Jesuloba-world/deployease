package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Jesuloba-world/deployease/backend/internal/infrastructure/database/dragonfly"
)

var (
	ErrNotFound             = redis.Nil
	ErrClientNotInitialized = errors.New("Dragonfly client not initialized")
)

type Store struct {
	client *redis.Client
}

func NewStore() (*Store, error) {
	client := dragonfly.GetClient()
	if client == nil {
		return nil, ErrClientNotInitialized
	}
	return &Store{
		client: client,
	}, nil
}

func (s *Store) Get(ctx context.Context, key string, value interface{}) error {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return err
	}
	return json.Unmarshal(data, value)
}

func (s *Store) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, expiration).Err()
}

func (s *Store) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *Store) HealthCheck(ctx context.Context) error {
	if s.client == nil {
		return ErrClientNotInitialized
	}
	return s.client.Ping(ctx).Err()
}
