package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/config"
)

type AIService struct {
	baseURL string
	client  *http.Client
}

type AnalyzeRequest struct {
	BlueprintID uuid.UUID `json:"blueprint_id"`
	S3Key       string    `json:"s3_key"`
}

type AnalyzeResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		baseURL: cfg.AI.ServiceURL,
		client: &http.Client{
			Timeout: cfg.AI.Timeout,
		},
	}
}

func (s *AIService) AnalyzeBlueprint(ctx context.Context, blueprintID uuid.UUID, s3Key string) (string, error) {
	reqBody := AnalyzeRequest{
		BlueprintID: blueprintID,
		S3Key:       s3Key,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/analyze", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	var result AnalyzeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		// Return raw response if not JSON
		return string(body), nil
	}

	if !result.Success {
		return "", fmt.Errorf("AI service error: %s", result.Error)
	}

	// Return the result as JSON string
	resultJSON, err := json.Marshal(result.Data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resultJSON), nil
}

func (s *AIService) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/health", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI service returned status %d", resp.StatusCode)
	}

	return nil
}
