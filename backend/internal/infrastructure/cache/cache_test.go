package cache

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jesuloba-world/deployease/backend/internal/config"
	"github.com/Jesuloba-world/deployease/backend/internal/infrastructure/database/dragonfly"
)

var testCacheStore *Store
var testContainer *dragonfly.DragonflyContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testContainer, err = dragonfly.StartDragonflyContainer(ctx)
	if err != nil {
		fmt.Printf("Error starting Dragonfly container: %v", err)
		os.Exit(1)
	}

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
	testCacheStore = tmpStore

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

func TestCacheStore_SetAndGet(t *testing.T) {
	ctx := context.Background()
	key := "cache:test:" + gonanoid.Must()
	type testStruct struct {
		Field1 string
		Field2 int
	}
	originalValue := testStruct{Field1: "hello", Field2: 123}
	var retrievedValue testStruct

	err := testCacheStore.Set(ctx, key, originalValue, 1*time.Hour)
	require.NoError(t, err, "failed to set cache item")

	err = testCacheStore.Get(ctx, key, &retrievedValue)
	require.NoError(t, err, "failed to get cache item")
	assert.Equal(t, originalValue, retrievedValue, "retrieved value does not match original")
}

func TestCacheStore_GetNonExistent(t *testing.T) {
	ctx := context.Background()
	key := "cache:non_existent:" + gonanoid.Must()
	var retrievedValue struct{ Field string }

	err := testCacheStore.Get(ctx, key, &retrievedValue)
	assert.ErrorIs(t, err, ErrNotFound, "expected ErrNotFound when key is non-existent")
}

func TestCacheStore_ItemExpiry(t *testing.T) {
	ctx := context.Background()
	key := "cache:expiry:" + gonanoid.Must()
	value := "expiring data"
	var retrievedValue string

	err := testCacheStore.Set(ctx, key, value, 100*time.Millisecond)
	require.NoError(t, err)

	time.Sleep(150 * time.Millisecond)

	err = testCacheStore.Get(ctx, key, &retrievedValue)
	assert.ErrorIs(t, err, ErrNotFound, "expected ErrNotFound after item expires")
}

func TestCacheStore_Delete(t *testing.T) {
	ctx := context.Background()
	key := "cache:delete:" + gonanoid.Must()
	value := "to be deleted"
	var retrievedValue string

	err := testCacheStore.Set(ctx, key, value, 1*time.Hour)
	require.NoError(t, err)

	// test if it is there
	err = testCacheStore.Get(ctx, key, &retrievedValue)
	require.NoError(t, err)
	assert.Equal(t, value, retrievedValue)

	// delete it now
	err = testCacheStore.Delete(ctx, key)
	assert.NoError(t, err, "failed to delete cache item")

	// vrify delete
	err = testCacheStore.Get(ctx, key, &retrievedValue)
	assert.ErrorIs(t, err, ErrNotFound, "expected ErrNotFound after deletion")
}

func TestCacheStore_HealthCheck(t *testing.T) {
	ctx := context.Background()
	// Assuming Dragonfly/Redis is running and accessible
	err := testCacheStore.HealthCheck(ctx)
	assert.NoError(t, err, "Health check should pass if Dragonfly is running")

	// Test health check when client is not initialized (simulate this by temporarily nil-ing it)
	// This is a bit of a hack for testing this specific error case.
	originalClient := testCacheStore.client
	testCacheStore.client = nil
	err = testCacheStore.HealthCheck(ctx)
	assert.ErrorIs(t, err, ErrClientNotInitialized, "Health check should fail if client is not initialized")
	testCacheStore.client = originalClient // Restore client
}

func TestNewStore_ClientNotInitialized(t *testing.T) {
	// Ensure dragonfly client is closed to simulate not initialized state
	currentClient := dragonfly.GetClient()
	if currentClient != nil {
		dragonfly.CloseClient()
	}

	_, err := NewStore()
	assert.ErrorIs(t, err, ErrClientNotInitialized, "NewStore should return ErrClientNotInitialized if dragonfly client is nil")

	// Re-initialize for other tests if it was open before
	if currentClient != nil {
		// Create a basic config for re-initialization
		cfg := &config.Config{
			Redis: config.RedisConfig{
				Host:     "localhost",
				Port:     "6379",
				Password: "",
				DB:       0,
			},
		}
		dragonfly.InitClient(cfg)
	}
}
