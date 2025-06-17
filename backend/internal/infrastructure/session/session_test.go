package session

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jesuloba-world/deployease/backend/internal/infrastructure/database/dragonfly"

)

var testSessionStore *Store
var testContainer *dragonfly.DragonflyContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testContainer, err = dragonfly.StartDragonflyContainer(ctx)
	if err != nil {
		fmt.Printf("Error starting Dragonfly container: %v", err)
		os.Exit(1)
	}

	fmt.Printf("test dragonfly container for session tests started at %s\n", testContainer.Address)

	cfg := testContainer.GetConfig()

	if err = dragonfly.InitClient(cfg); err != nil {
		fmt.Printf("Error initializing Dragonfly client: %v", err)
		testContainer.Cleanup(ctx)
		os.Exit(1)
	}

	client := dragonfly.GetClient()
	if client == nil {
		fmt.Printf("Error getting Dragonfly client: %v", err)
		testContainer.Cleanup(ctx)
		os.Exit(1)
	}

	tmpStore, err := NewStore()
	if err != nil {
		fmt.Printf("Error initializing session store: %v", err)
		testContainer.Cleanup(ctx)
		os.Exit(1)
	}
	testSessionStore = tmpStore

	keys, err := client.Keys(ctx, "session:*").Result()
	if err != nil {
		fmt.Printf("Warning: failed to scan for existing session keys for cleanup: %v\n", err)
	} else if len(keys) > 0 {
		if err := client.Del(ctx, keys...).Err(); err != nil {
			fmt.Printf("Warning: Failed to delete existing session keys: %v\n", err)
		}
	}

	code := m.Run()

	dragonfly.CloseClient()
	if err := testContainer.Cleanup(ctx); err != nil {
		fmt.Printf("Could not terminate test container: %v", err)
	}
	os.Exit(code)
}

func TestSessionStore_SetAndGetSession(t *testing.T) {
	ctx := context.Background()
	sessionID := gonanoid.Must()
	userID := gonanoid.Must()
	expiresAt := time.Now().Add(1 * time.Hour)

	sess := &Session{
		ID:        sessionID,
		UserID:    userID,
		Data:      map[string]interface{}{"key1": "value1", "key2": 123},
		ExpiresAt: expiresAt,
	}

	// Expected data after JSON round-trip (numbers become float64)
	expectedData := map[string]interface{}{"key1": "value1", "key2": float64(123)}

	err := testSessionStore.Set(ctx, sess)
	require.NoError(t, err, "failed to set session")

	retrievedSess, err := testSessionStore.Get(ctx, sessionID)
	require.NoError(t, err, "failed to get session")
	require.NotNil(t, retrievedSess, "expected session to be not nil")

	assert.Equal(t, sess.UserID, retrievedSess.UserID, "expected user ID to match")
	assert.Equal(t, sess.ID, retrievedSess.ID, "expected session ID to match")
	assert.Equal(t, expectedData, retrievedSess.Data, "expected data to match")

	assert.WithinDuration(t, sess.ExpiresAt, retrievedSess.ExpiresAt, time.Second)
}

func TestSessionStore_GetNonExistentSession(t *testing.T) {
	ctx := context.Background()
	nonExistentID := gonanoid.Must()

	retrievedSess, err := testSessionStore.Get(ctx, nonExistentID)
	require.NoError(t, err, "expected no error for non-existent session")
	assert.Nil(t, retrievedSess, "expected session to be nil")
}

func TestSessionStore_SessionExpiry(t *testing.T) {
	ctx := context.Background()
	sessionID := gonanoid.Must()
	expiresAt := time.Now().Add(100 * time.Millisecond) // short expiry for testing

	sess := &Session{
		ID:        sessionID,
		UserID:    gonanoid.Must(),
		ExpiresAt: expiresAt,
	}

	err := testSessionStore.Set(ctx, sess)
	require.NoError(t, err)

	// wait for session to expire
	time.Sleep(150 * time.Millisecond)

	retrievedSess, err := testSessionStore.Get(ctx, sessionID)
	assert.NoError(t, err, "expected no error for expired session")
	assert.Nil(t, retrievedSess, "expected session to be nil after expiry")

	val, err := dragonfly.GetClient().Get(ctx, "session:"+sessionID).Result()
	assert.ErrorIs(t, err, redis.Nil, "expected redis.Nil error for expired session key")
	assert.Empty(t, val, "expected session key to be empty after expiry")
}

func TestSessionStore_DeleteSession(t *testing.T) {
	ctx := context.Background()
	sessionID := gonanoid.Must()
	expiresAt := time.Now().Add(1 * time.Hour)

	sess := &Session{
		ID:        sessionID,
		UserID:    gonanoid.Must(),
		ExpiresAt: expiresAt,
	}

	err := testSessionStore.Set(ctx, sess)
	require.NoError(t, err)

	retrievedSess, err := testSessionStore.Get(ctx, sessionID)
	require.NoError(t, err)
	require.NotNil(t, retrievedSess, "expected session to be not nil after retrieval")

	err = testSessionStore.Delete(ctx, sessionID)
	assert.NoError(t, err, "failed to delete session")

	deletedSess, err := testSessionStore.Get(ctx, sessionID)
	assert.NoError(t, err)
	assert.Nil(t, deletedSess, "expected session to be nil after deletion")
}

func TestSessionStore_SetSessionWithPastExpiry(t *testing.T) {
	ctx := context.Background()
	sessionID := gonanoid.Must()
	expiresAt := time.Now().Add(-1 * time.Hour) // expired session

	sess := &Session{
		ID:        sessionID,
		UserID:    gonanoid.Must(),
		ExpiresAt: expiresAt,
	}

	err := testSessionStore.Set(ctx, sess)
	assert.Error(t, err, "expected error for past expiry")
}
