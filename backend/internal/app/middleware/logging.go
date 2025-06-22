package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/uptrace/bunrouter"

)

type LoggingConfig struct {
	SkipPaths []string
}

func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		SkipPaths: []string{"/health", "/ready"},
	}
}

func Logging(config LoggingConfig) bunrouter.MiddlewareFunc {

	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			start := time.Now()

			for _, skipPath := range config.SkipPaths {
				if req.URL.Path == skipPath {
					return next(w, req)
				}
			}

			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			err := next(wrapper, req)

			duration := time.Since(start)

			log.Printf(
				"%s %s %d %v %s",
				req.Method,
				req.URL.Path,
				wrapper.statusCode,
				duration,
				req.RemoteAddr,
			)

			return err
		}
	}
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(data)
}
