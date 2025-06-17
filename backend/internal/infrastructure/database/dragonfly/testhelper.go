package dragonfly

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Jesuloba-world/deployease/backend/internal/config"
)

type DragonflyContainer struct {
	Container testcontainers.Container
	Host      string
	Port      string
	Address   string
}

func StartDragonflyContainer(ctx context.Context) (*DragonflyContainer, error) {
	redisContainer, err := redis.Run(ctx,
		"docker.dragonflydb.io/dragonflydb/dragonfly:latest",
		testcontainers.WithWaitStrategy(wait.ForLog("AcceptServer - listening on port 6379")),
	)
	if err != nil {
		return nil, fmt.Errorf("could not start Dragonfly container: %w", err)
	}

	host, err := redisContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get Dragonfly host: %w", err)
	}

	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, fmt.Errorf("could not get Dragonfly port: %w", err)
	}

	address := fmt.Sprintf("%s:%s", host, port.Port())

	return &DragonflyContainer{
		Container: redisContainer,
		Host:      host,
		Port:      port.Port(),
		Address:   address,
	}, nil
}

func (dc *DragonflyContainer) GetConfig() *config.Config {
	return &config.Config{
		Redis: config.RedisConfig{
			Host:     dc.Host,
			Port:     dc.Port,
			Password: "",
			DB:       0,
		},
	}
}

func (dc *DragonflyContainer) Cleanup(ctx context.Context) error {
	if dc.Container != nil {
		return dc.Container.Terminate(ctx)
	}
	return nil
}
