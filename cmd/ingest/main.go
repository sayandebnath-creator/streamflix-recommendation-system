package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"streamflix-backend/internal/config"
	"streamflix-backend/internal/database"
	"streamflix-backend/internal/embedding"
	"streamflix-backend/internal/ingestion"
	"streamflix-backend/internal/movie"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	slog.SetDefault(logger)

	cfg := config.Load()

	database.Connect(cfg)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	embeddingService := embedding.NewHTTPService(
		httpClient,
		cfg.EmbeddingServiceURL,
	)

	movieRepository := movie.NewRepository(database.DB)

	ingestionService := ingestion.NewService(
		movieRepository,
		embeddingService,
		cfg.EmbeddingWorkers,
	)

	slog.Info("Starting embedding ingestion...")

	if err := ingestionService.ProcessPendingEmbeddings(
		context.Background(),
		cfg.EmbeddingBatchSize,
	); err != nil {
		slog.Error("embedding ingestion failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Embedding ingestion completed.")
}
