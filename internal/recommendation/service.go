package recommendation

import (
	"context"

	"streamflix-backend/internal/movie"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) SimilarMovies(
	ctx context.Context,
	id string,
	limit int,
) ([]movie.Movie, error) {
	return s.repository.FindSimilarMovies(ctx, id, limit)
}