package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Jesuloba-world/deployease/backend/internal/infrastructure/database/dragonfly"
)

var (
	ErrClientNotInitialized = errors.New("Dragonfly client not initialized")
)

type Session struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
	ExpiresAt time.Time              `json:"expires_at"`
}

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

func (s *Store) Get(ctx context.Context, sessionID string) (*Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var sess Session
	err = json.Unmarshal([]byte(val), &sess)
	if err != nil {
		return nil, err
	}

	if sess.ExpiresAt.Before(time.Now()) {
		_ = s.Delete(ctx, sessionID)
		return nil, nil
	}

	return &sess, nil
}

func (s *Store) Set(ctx context.Context, session *Session) error {
	key := fmt.Sprintf("session:%s", session.ID)
	val, err := json.Marshal(session)
	if err != nil {
		return err
	}

	duration := time.Until(session.ExpiresAt)
	if duration <= 0 {
		return errors.New("session expiry must be some time in the future")
	}

	return s.client.Set(ctx, key, val, duration).Err()
}

func (s *Store) Delete(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.client.Del(ctx, key).Err()
}
