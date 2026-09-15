package movie


import (
	"context"

	"github.com/pgvector/pgvector-go"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAllMovies() ([]Movie, error) {
	return s.repo.GetAll()
}

func (s *Service) GetMovies(
	ctx context.Context,
	page int,
	limit int,
) ([]Movie, int64, error) {
	return s.repo.GetPaginated(
		ctx,
		page,
		limit,
	)
}

func (s *Service) CreateMovie(movie *Movie) error {
	return s.repo.Create(movie)
}

func (s *Service) SearchSimilarMovies(
	ctx context.Context,
	vector pgvector.Vector,
	limit int,
) ([]Movie, error) {
	return s.repo.SearchSimilarMovies(
		ctx,
		vector,
		limit,
	)
}

func (s *Service) GetMovie(
    ctx context.Context,
    id uuid.UUID,
) (*Movie, error) {
    return s.repo.GetByID(ctx, id)
}