package middleware

import (
	"context"
	"net/http"
	"regexp"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/uptrace/bunrouter"
)

type RequestIDGenerator func() string

type RequestIDValidator func(string) bool

type RequestIDConfig struct {
	HeaderName    string
	ContextKey    string
	Generator     RequestIDGenerator
	Validator     RequestIDValidator
	ReuseExisting bool
}

func DefaultRequestIDConfig() RequestIDConfig {
	return RequestIDConfig{
		HeaderName:    "X-Request-ID",
		ContextKey:    "request_id",
		Generator:     DefaultRequestIDGenerator,
		Validator:     DefaultRequestIDValidator,
		ReuseExisting: true,
	}
}

func DefaultRequestIDGenerator() string {
	return gonanoid.Must()
}

func DefaultRequestIDValidator(id string) bool {
	if len(id) != 21 {
		return false
	}
	// Nano ID uses URL-safe characters: A-Za-z0-9_-
	matched, _ := regexp.MatchString(`^[A-Za-z0-9_-]{21}$`, id)
	return matched
}

func RequestID(config RequestIDConfig) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			var requestID string

			if config.ReuseExisting {
				existingID := req.Header.Get(config.HeaderName)
				if existingID != "" && config.Validator(existingID) {
					requestID = existingID
				}
			}

			if requestID == "" {
				requestID = config.Generator()
			}

			w.Header().Set(config.HeaderName, requestID)

			ctx := req.Context()
			ctx = context.WithValue(ctx, config.ContextKey, requestID)
			req = req.WithContext(ctx)

			return next(w, req)
		}
	}
}
