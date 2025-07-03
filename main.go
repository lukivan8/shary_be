package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"shary_be/internal/config"
	"shary_be/internal/handlers"
	"shary_be/internal/repository"
	"shary_be/internal/router"
	"shary_be/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger based on configuration
	var logger *zap.Logger
	var err error

	switch cfg.LogLevel {
	case "debug":
		logger, err = zap.NewDevelopment()
	case "production":
		logger, err = zap.NewProduction()
	default:
		// Default to info level
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		logger, err = config.Build()
	}

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting application",
		zap.String("environment", cfg.Environment),
		zap.String("log_level", cfg.LogLevel),
		zap.Int("port", cfg.Port),
	)

	// Initialize database connection
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("Successfully connected to database")

	// Initialize repositories
	itemRepo := repository.NewItemRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize services
	itemService := service.NewItemService(itemRepo, logger)
	categoryService := service.NewCategoryService(categoryRepo, logger)

	// Initialize handlers
	itemHandler := handlers.NewItemHandler(itemService, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)

	// Setup Chi router
	handler := router.SetupRouter(itemHandler, categoryHandler, logger)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", zap.Int("port", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
