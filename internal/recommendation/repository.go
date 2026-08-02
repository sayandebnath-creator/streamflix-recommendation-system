package recommendation

import (
	"context"

	"streamflix-backend/internal/movie"
)

type Repository interface {
	FindSimilarMovies(
		ctx context.Context,
		id string,
		limit int,
	) ([]movie.Movie, error)
}