package embedding

import "context"

type Service interface {
	GenerateEmbedding(ctx context.Context, text string, isQuery bool) ([]float32, error)
}