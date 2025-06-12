package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Jesuloba-world/deployease/backend/internal/config"
)

type TestContainer struct {
	Container *postgres.PostgresContainer
	Pool      *pgxpool.Pool
	DSN       string
}

func SetupTestContainer(t *testing.T) *TestContainer {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start postgres container: %v", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get postgres container host: %v", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get postgres container port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable&pool_max_conns=10&pool_min_conns=2&pool_max_conn_lifetime=1m&pool_max_conn_idle_time=30s",
		host, port.Port())

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed to create pgxpool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping database with pgx pool: %v", err)
	}

	return &TestContainer{
		Container: postgresContainer,
		Pool:      pool,
		DSN:       dsn,
	}
}

func (tc *TestContainer) Cleanup(t *testing.T) {
	ctx := context.Background()

	if tc.Pool != nil {
		tc.Pool.Close()
	}

	if tc.Container != nil {
		if err := tc.Container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate postgres container: %v", err)
		}
	}
}

func SetupTestManager(t *testing.T) (*Manager, func()) {
	tc := SetupTestContainer(t)

	// Create config from test container
	port, err := tc.Container.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}
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
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	cleanup := func() {
		manager.Close()
		tc.Cleanup(t)
	}

	return manager, cleanup
}

func setupTestDB(t *testing.T) *pgxpool.Pool {
	tc := SetupTestContainer(t)
	t.Cleanup(func() {
		tc.Cleanup(t)
	})
	return tc.Pool
}
