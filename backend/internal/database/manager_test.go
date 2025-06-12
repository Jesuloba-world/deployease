package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Jesuloba-world/deployease/backend/internal/config"
)

func TestNewManager(t *testing.T) {
	tc := SetupTestContainer(t)
	defer tc.Cleanup(t)

	port, err := tc.Container.MappedPort(context.Background(), "5432")
	require.NoError(t, err)

	cfg := &config.DatabaseConfig{
		Host:            "localhost",
		Port:            port.Port(),
		User:            "testuser",
		Password:        "testpass",
		DBName:          "testdb",
		SSLMode:         "disable",
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: 1 * time.Minute,
		ConnMaxIdleTime: 30 * time.Second,
	}

	manager, err := NewManager(cfg)
	require.NoError(t, err)
	defer manager.Close()
	require.NotNil(t, manager)
}

func TestNewManagerWithNilCOnfig(t *testing.T) {
	manager, err := NewManager(nil)
	assert.Error(t, err)
	assert.Nil(t, manager)
}

func TestManagerClose(t *testing.T) {
	tc := SetupTestContainer(t)
	defer tc.Cleanup(t)

	port, err := tc.Container.MappedPort(context.Background(), "5432")
	require.NoError(t, err)

	cfg := &config.DatabaseConfig{
		Host:            "localhost",
		Port:            port.Port(),
		User:            "testuser",
		Password:        "testpass",
		DBName:          "testdb",
		SSLMode:         "disable",
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: 1 * time.Minute,
		ConnMaxIdleTime: 30 * time.Second,
	}

	manager, err := NewManager(cfg)
	require.NoError(t, err)

	// Test close
	manager.Close()
}
