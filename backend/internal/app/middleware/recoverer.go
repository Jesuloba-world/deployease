package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/uptrace/bunrouter"
)

type RecovererConfig struct {
	EnableStackTrace   bool
	EnablePanicLogs    bool
	CustomErrorHandler func(w http.ResponseWriter, req bunrouter.Request, err interface{})
}

func DefaultRecovererConfig() RecovererConfig {
	return RecovererConfig{
		EnableStackTrace:   true,
		EnablePanicLogs:    true,
		CustomErrorHandler: nil,
	}
}

func Recoverer(config RecovererConfig) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) (err error) {
			defer func() {
				if r := recover(); r != nil {
					// log the panic if panic log is enabled
					if config.EnablePanicLogs {
						log.Printf("Panic recovered: %v", r)
					}
					// log stack trace if enabled
					if config.EnableStackTrace {
						log.Printf("Stack trace: %s", debug.Stack())
					}

					// use customErrHandler if provided
					if config.CustomErrorHandler != nil {
						config.CustomErrorHandler(w, req, r)
						err = nil
						return
					}

					// default error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					errorResponse := fmt.Sprintf(`{"error":"Internal server error","message":"An unexpected error occurred","status":%d}`, http.StatusInternalServerError)
					w.Write([]byte(errorResponse))

					// stop panic
					err = nil
				}
			}()

			return next(w, req)
		}
	}
}

func SimpleRecoverer() bunrouter.MiddlewareFunc {
	return Recoverer(DefaultRecovererConfig())
}

func RecovererWithCustomHandler(handler func(w http.ResponseWriter, req bunrouter.Request, err interface{})) bunrouter.MiddlewareFunc {
	config := DefaultRecovererConfig()
	config.CustomErrorHandler = handler
	return Recoverer(config)
}
