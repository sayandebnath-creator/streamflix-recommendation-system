package ingestion

import (
	"context"
	"log/slog"
	"sync"

	"streamflix-backend/internal/movie"
)

type Worker struct {
	service *Service
}

func NewWorker(service *Service) *Worker {
	return &Worker{
		service: service,
	}
}

func (w *Worker) Start(
	ctx context.Context,
	id int,
	jobs <-chan movie.Movie,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case m, ok := <-jobs:
			if !ok {
				return
			}

			if err := w.service.processMovie(ctx, m); err != nil {
				slog.Error(
					"worker failed processing movie",
					"worker_id",
					id,
					"movie_id",
					m.ID,
					"error",
					err,
				)
			}
		}
	}
}