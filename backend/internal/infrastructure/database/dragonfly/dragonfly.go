package dragonfly

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Jesuloba-world/deployease/backend/internal/config"
)

var rdb *redis.Client

func InitClient(cfg *config.Config) error {
	host := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	password := cfg.Redis.Password
	db := cfg.Redis.DB

	fmt.Printf("Connecting to Dragonfly at %s\n", host)

	rdb = redis.NewClient(&redis.Options{
		Addr:         host,
		Password:     password,
		DB:           db,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Dragonfly: %w", err)
	}

	fmt.Println("Successfully connected to Dragonfly!")
	return nil
}

func GetClient() *redis.Client {
	return rdb
}

func CloseClient() error {
	if rdb != nil {
		err := rdb.Close()
		if err != nil {
			return fmt.Errorf("failed to close Dragonfly connection: %w", err)
		}
		rdb = nil
	}
	return nil
}
