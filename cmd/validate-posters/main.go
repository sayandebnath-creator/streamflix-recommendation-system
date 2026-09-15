package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"streamflix-backend/internal/config"
	"streamflix-backend/internal/database"
	"streamflix-backend/internal/movie"
)

const tmdbImageBaseURL = "https://image.tmdb.org/t/p/w500"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()
	database.Connect(cfg)

	movieRepository := movie.NewRepository(database.DB)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	for {
		movies, err := movieRepository.GetMoviesForPosterValidation(
			context.Background(),
			cfg.PosterValidationBatchSize,
		)
		if err != nil {
			slog.Error(
				"failed to fetch movies for poster validation",
				"error", err,
			)
			os.Exit(1)
		}

		if len(movies) == 0 {
			break
		}

		slog.Info(
			"starting poster validation batch",
			"movies", len(movies),
			"workers", cfg.PosterValidationWorkers,
		)

		jobs := make(chan movie.Movie)
		var wg sync.WaitGroup

		for i := 0; i < cfg.PosterValidationWorkers; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for m := range jobs {
					valid := validatePoster(client, m.PosterPath)

					err := movieRepository.UpdatePosterValidity(
						context.Background(),
						m.ID,
						valid,
					)
					if err != nil {
						slog.Error(
							"failed to update poster validity",
							"movie_id", m.ID,
							"title", m.Title,
							"error", err,
						)
						continue
					}
				}
			}()
		}

		for _, m := range movies {
			jobs <- m
		}

		close(jobs)
		wg.Wait()

		slog.Info(
			"poster validation batch completed",
			"movies", len(movies),
		)
	}

	slog.Info("poster validation completed")
}

func validatePoster(client *http.Client, posterPath string) bool {
	if posterPath == "" {
		return false
	}

	url := tmdbImageBaseURL + posterPath

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)

	return resp.StatusCode == http.StatusOK
}