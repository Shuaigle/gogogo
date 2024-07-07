package main

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (app *application) logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        defer func() {
            app.logger.Info("Request",
                zap.String("method", r.Method),
                zap.String("url", r.URL.String()),
                zap.Duration("duration", time.Since(start)))
        }()
        next.ServeHTTP(w, r)
    })
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
					if err := recover(); err != nil {
							w.Header().Set("Connection", "close")
							app.serverError(w, fmt.Errorf("%s", err))
					}
			}()
			next.ServeHTTP(w, r)
	})
}
