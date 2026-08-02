package ingestion

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	
	"sync"
	"streamflix-backend/internal/embedding"
	"streamflix-backend/internal/movie"
)

type MovieRepository interface {
	GetWithoutEmbeddings(ctx context.Context, limit int) ([]movie.Movie, error)
	UpdateEmbedding(ctx context.Context, movieID uuid.UUID, embedding []float32) error
}
type Service struct {
	movieRepository  MovieRepository
	embeddingService embedding.Service
	workers          int
}

func NewService(
	movieRepository MovieRepository,
	embeddingService embedding.Service,
	workers int,
) *Service {
	return &Service{
		movieRepository:  movieRepository,
		embeddingService: embeddingService,
		workers:          workers,
	}
}

func (s *Service) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	return s.embeddingService.GenerateEmbedding(ctx, text)
}

func (s *Service) processMovie(ctx context.Context, movie movie.Movie) error {
	if movie.Overview == "" {
		slog.Info(
			"skipping movie: empty overview",
			"movie_id", movie.ID,
			"title", movie.Title,
		)
		return nil
	}

	slog.Info(
		"generating embedding",
		"movie_id", movie.ID,
		"title", movie.Title,
	)

	vector, err := s.embeddingService.GenerateEmbedding(ctx, movie.Overview)
	if err != nil {
		return err
	}

	if err := s.movieRepository.UpdateEmbedding(ctx, movie.ID, vector); err != nil {
		return err
	}

	slog.Info(
		"embedding saved",
		"movie_id", movie.ID,
		"title", movie.Title,
	)

	return nil
}

func (s *Service) ProcessPendingEmbeddings(ctx context.Context, limit int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	for {
		movies, err := s.movieRepository.GetWithoutEmbeddings(ctx, limit)
		if err != nil {
			return err
		}
		if len(movies) == 0 {
			break
		}
		jobs := make(chan movie.Movie, limit)

		var wg sync.WaitGroup

		worker := NewWorker(s)

		for i := 0; i < s.workers; i++ {
			wg.Add(1)

			go worker.Start(
				ctx,
				i,
				jobs,
				&wg,
			)
		}

		for _, movie := range movies {
			jobs <- movie
		}

		close(jobs)

		wg.Wait()
	}

	return nil
}
