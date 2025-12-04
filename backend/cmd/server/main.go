package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/config"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/handlers"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/middleware"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/repository"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/services"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting Construction Estimation & Bidding API",
		"version", "1.0.0",
		"env", cfg.Server.Env,
		"port", cfg.Server.Port)

	// Initialize database
	db, err := repository.NewDatabase(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	projectRepo := repository.NewProjectRepository(db)
	blueprintRepo := repository.NewBlueprintRepository(db)
	jobRepo := repository.NewJobRepository(db)
	bidRepo := repository.NewBidRepository(db)

	// Initialize services
	s3Service, err := services.NewS3Service(cfg)
	if err != nil {
		slog.Error("Failed to initialize S3 service", "error", err)
		os.Exit(1)
	}

	// Ensure S3 bucket exists
	if err := s3Service.EnsureBucket(context.Background()); err != nil {
		slog.Warn("Failed to ensure S3 bucket exists", "error", err)
		// Don't exit - bucket might exist already or will be created by admin
	}

	aiService := services.NewAIService(cfg)

	// Initialize worker
	worker := services.NewWorker(jobRepo, blueprintRepo, aiService, cfg)
	ctx, cancel := context.WithCancel(context.Background())
	worker.Start(ctx)
	defer func() {
		cancel()
		worker.Stop()
	}()

	// Initialize handlers
	handler := handlers.NewHandler(db, projectRepo, blueprintRepo, jobRepo, bidRepo, s3Service, aiService)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS)

	// Routes
	r.Get("/", handler.Root)
	r.Get("/health", handler.Health)

	// Blueprint upload routes
	r.Post("/projects/{id}/blueprints/upload-url", handler.CreateUploadURL)
	r.Post("/blueprints/{id}/complete-upload", handler.CompleteUpload)

	// Blueprint analysis routes
	r.Get("/blueprints/{id}/analysis", handler.GetBlueprintAnalysis)
	r.Get("/blueprints/{id}/takeoff-summary", handler.GetBlueprintTakeoffSummary)

	// Job routes
	r.Post("/blueprints/{id}/analyze", handler.AnalyzeBlueprint)
	r.Get("/jobs/{id}", handler.GetJobStatus)

	// Bid routes
	r.Get("/projects/{id}/pricing-summary", handler.GetPricingSummary)
	r.Post("/projects/{id}/generate-bid", handler.GenerateBid)
	r.Get("/projects/{id}/bids", handler.GetProjectBids)
	r.Get("/bids/{id}", handler.GetBid)
	r.Get("/bids/{id}/pdf", handler.GetBidPDF)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("Server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}
