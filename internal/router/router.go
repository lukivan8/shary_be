package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"shary_be/internal/handlers"
)

// SetupRouter creates and configures the Chi router with all routes and middleware
func SetupRouter(itemHandler *handlers.ItemHandler, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	// Add Chi's built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Add custom middleware
	r.Use(loggingMiddleware(logger))
	r.Use(recoveryMiddleware(logger))

	// Health check endpoint
	r.Get("/health", healthCheckHandler)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/items", func(r chi.Router) {
			r.Get("/", itemHandler.GetAllItems)
			r.Post("/", itemHandler.CreateItem)
			r.Get("/location/{location}", itemHandler.GetItemsByLocation)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", itemHandler.GetItemByID)
				r.Put("/", itemHandler.UpdateItem)
				r.Delete("/", itemHandler.DeleteItem)
			})
		})
	})

	return r
}

// healthCheckHandler handles the health check endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

// recoveryMiddleware recovers from panics
func recoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						zap.Any("error", err),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
					)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
