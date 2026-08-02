package recommendation

import (
	"context"

	"gorm.io/gorm"

	"streamflix-backend/internal/movie"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) FindSimilarMovies(
	ctx context.Context,
	id string,
	limit int,
) ([]movie.Movie, error) {
	var movies []movie.Movie

	query := `
SELECT *
FROM movies
WHERE id <> ?
  AND embedding IS NOT NULL
ORDER BY embedding <=> (
	SELECT embedding
	FROM movies
	WHERE id = ?
)
LIMIT ?;
`

	err := r.db.WithContext(ctx).
		Raw(query, id, id, limit).
		Scan(&movies).Error
	if err != nil {
		return nil, err
	}

	return movies, nil
}