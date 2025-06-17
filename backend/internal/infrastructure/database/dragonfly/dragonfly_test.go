package dragonfly

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testContainer *DragonflyContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testContainer, err = StartDragonflyContainer(ctx)
	if err != nil {
		fmt.Printf("Could not start dragonfly container: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Test Dragonfly container started at %s\n", testContainer.Address)

	code := m.Run()

	if err := testContainer.Cleanup(ctx); err != nil {
		fmt.Printf("Could not cleanup dragonfly container: %s\n", err)
	}

	os.Exit(code)
}

func TestInitClient_Success(t *testing.T) {
	if testContainer == nil {
		t.Skip("test container not available, skipping test")
	}

	cfg := testContainer.GetConfig()

	err := InitClient(cfg)
	require.NoError(t, err, "Initclient should succeed")

	client := GetClient()
	require.NotNil(t, client, "client should be initialized")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	assert.NoError(t, err, "Ping should succeed with running dragonfly server")

	CloseClient()
	assert.Nil(t, GetClient(), "client should be cleaned up")
}

func TestInitClient_WithOptions(t *testing.T) {
	if testContainer == nil {
		t.Skip("test container not available, skipping test")
	}

	cfg := testContainer.GetConfig()
	cfg.Redis.DB = 1

	err := InitClient(cfg)
	require.NoError(t, err, "InitClient should succeed with custom config")

	client := GetClient()
	require.NotNil(t, client, "client should be initialized")
	assert.Equal(t, 1, client.Options().DB, "DB should be set to custom value")
	CloseClient()
}

func TestGetClient_Uninitialized(t *testing.T) {
	if GetClient() != nil {
		CloseClient()
	}
	assert.Nil(t, GetClient(), "client should be cleaned up")
}

func TestCloseClient_Idempotent(t *testing.T) {
	if testContainer == nil {
		t.Skip("test container not available, skipping test")
	}

	cfg := testContainer.GetConfig()

	err := InitClient(cfg)
	require.NoError(t, err)
	require.NotNil(t, GetClient(), "client should be initialized")
	CloseClient()
	assert.Nil(t, GetClient(), "client should be cleaned up")
	assert.NotPanics(t, func() {
		CloseClient()
	}, "CloseClient should be idempotent")
	assert.Nil(t, GetClient(), "client should remain nil after second call")
}
