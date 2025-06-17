package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Jesuloba-world/deployease/backend/internal/config"
)

type Manager struct {
	dbPool *pgxpool.Pool
}

func NewManager(cfg *config.DatabaseConfig) (*Manager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database configuration is required")
	}

	dsn := cfg.GetDSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN: %w", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Manager{dbPool: dbPool}, nil
}

func (m *Manager) Close() {
	if m.dbPool != nil {
		m.dbPool.Close()
	}
}

func (m *Manager) DBPool() *pgxpool.Pool {
	return m.dbPool
}
