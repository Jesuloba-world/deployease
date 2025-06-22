package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/uptrace/bunrouter"

	"github.com/Jesuloba-world/deployease/backend/internal/config"
)

type TimeoutConfig struct {
	Timeout time.Duration
}

func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Timeout: 30 * time.Second,
	}
}

func NewTimeoutConfigFromServer(cfg *config.Config) TimeoutConfig {
	timeout := 30 * time.Second
	if cfg.Server.ReadTimeout > 0 {
		timeout = cfg.Server.ReadTimeout
	}

	return TimeoutConfig{
		Timeout: timeout,
	}
}

func Timeout(config TimeoutConfig) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			ctx := req.Context()
			ctx, cancel := context.WithTimeout(ctx, config.Timeout)
			defer cancel()

			req = req.WithContext(ctx)

			return next(w, req)
		}
	}
}
