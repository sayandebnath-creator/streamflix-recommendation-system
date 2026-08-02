package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPService struct {
	client      *http.Client
	baseURL     string
	maxRetries  int
	baseBackoff time.Duration
}

type embeddingRequest struct {
	Text string `json:"text"`
}

type embeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func shouldRetry(statusCode int) bool {
	return statusCode >= 500 || statusCode == 0
}

func NewHTTPService(client *http.Client, baseURL string) *HTTPService {
	return &HTTPService{
		client:      client,
		baseURL:     baseURL,
		maxRetries:  3,
		baseBackoff: time.Second,
	}
}

func (s *HTTPService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	request := embeddingRequest{
		Text: text,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal embedding request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL+"/embed",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create embedding request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response

	for attempt := 0; attempt <= s.maxRetries; attempt++ {

		resp, err = s.client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		if err == nil && !shouldRetry(resp.StatusCode) {
			return nil, fmt.Errorf(
				"embedding service returned status %d",
				resp.StatusCode,
			)
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempt == s.maxRetries {
			return nil, fmt.Errorf("embedding request failed after %d retries", s.maxRetries)
		}

		delay := s.baseBackoff * time.Duration(1<<attempt)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"embedding service returned status %d",
			resp.StatusCode,
		)
	}

	var response embeddingResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode embedding response: %w", err)
	}

	return response.Embedding, nil
}
