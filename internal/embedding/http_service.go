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
	Text    string `json:"text"`
	IsQuery bool   `json:"is_query"`
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

func (s *HTTPService) GenerateEmbedding(
	ctx context.Context,
	text string,
	isQuery bool,
) ([]float32, error) {
	requestBody := embeddingRequest{
		Text:    text,
		IsQuery: isQuery,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal embedding request: %w", err)
	}

	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			s.baseURL+"/embed",
			bytes.NewReader(body),
		)
		if err != nil {
			return nil, fmt.Errorf("create embedding request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := s.client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()

			var response embeddingResponse

			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				return nil, fmt.Errorf(
					"decode embedding response: %w",
					err,
				)
			}

			return response.Embedding, nil
		}

		if err == nil && !shouldRetry(resp.StatusCode) {
			status := resp.StatusCode
			resp.Body.Close()

			return nil, fmt.Errorf(
				"embedding service returned status %d",
				status,
			)
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempt == s.maxRetries {
			if err != nil {
				return nil, fmt.Errorf(
					"embedding request failed after %d retries: %w",
					s.maxRetries,
					err,
				)
			}

			return nil, fmt.Errorf(
				"embedding service returned status %d after %d retries",
				resp.StatusCode,
				s.maxRetries,
			)
		}

		delay := s.baseBackoff * time.Duration(1<<attempt)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("embedding request failed")
}