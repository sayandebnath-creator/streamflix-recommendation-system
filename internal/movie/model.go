package movie

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type Movie struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// TMDB Metadata
	TMDBID int `gorm:"uniqueIndex" json:"tmdb_id"`

	Title       string    `gorm:"not null" json:"title"`
	Overview    string    `gorm:"type:text" json:"overview"`
	Genres      string    `gorm:"type:text" json:"genres"`
	Language    string    `gorm:"size:10" json:"language"`
	ReleaseDate *time.Time `json:"release_date"`

	Runtime int `json:"runtime"`

	Popularity float64 `json:"popularity"`

	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`

	PosterPath string `json:"poster_path"`

	// Semantic embedding generated using BAAI/bge-small-en-v1.5
	Embedding pgvector.Vector `gorm:"type:vector(384)" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Movie) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}