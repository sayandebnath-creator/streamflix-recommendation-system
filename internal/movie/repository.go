package movie

import (
	"context"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll() ([]Movie, error) {
	var movies []Movie

	err := r.db.Find(&movies).Error
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *Repository) GetWithoutEmbeddings(ctx context.Context, limit int) ([]Movie, error) {
	var movies []Movie

	err := r.db.WithContext(ctx).
		Where("embedding IS NULL").
		Limit(limit).
		Find(&movies).Error

	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *Repository) UpdateEmbedding(ctx context.Context, movieID uuid.UUID, embedding []float32) error {
	vector := pgvector.NewVector(embedding)
	return r.db.WithContext(ctx).
		Model(&Movie{}).
		Where("id = ?", movieID).
		Update("embedding", vector).
		Error
}

func (r *Repository) Create(movie *Movie) error {
	return r.db.Create(movie).Error
}

func (r *Repository) CreateBatch(movies []Movie) error {
	return r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tmdb_id"}},
			DoNothing: true,
		}).
		CreateInBatches(movies, 1000).Error
}


func (r *Repository) SearchSimilarMovies(
	ctx context.Context,
	vector pgvector.Vector,
	limit int,
) ([]Movie, error) {
	var movies []Movie

	err := r.db.WithContext(ctx).
		Where("embedding IS NOT NULL").
		Order(
			gorm.Expr(
				"embedding <=> ?",
				vector,
			),
		).
		Limit(limit).
		Find(&movies).Error

	if err != nil {
		return nil, err
	}

	return movies, nil
}