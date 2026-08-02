package embedding

import "context"

type Service interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}