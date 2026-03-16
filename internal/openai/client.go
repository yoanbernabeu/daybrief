package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	apiKey          string
	model           string
	language        string
	maxHighlights   int
	editorialPrompt string
	logger          *slog.Logger
	httpClient      *http.Client
}

func NewClient(apiKey, model, language string, maxHighlights int, editorialPrompt string, logger *slog.Logger) *Client {
	return &Client{
		apiKey:          apiKey,
		model:           model,
		language:        language,
		maxHighlights:   maxHighlights,
		editorialPrompt: editorialPrompt,
		logger:          logger,
		httpClient:      &http.Client{Timeout: 60 * time.Second},
	}
}

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *Client) completeJSON(ctx context.Context, prompt string) (string, error) {
	body := chatCompletionRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: "You are a precise assistant. Return only valid JSON, with no markdown and no extra text."},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.2,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshaling openai request: %w", err)
	}

	result, err := withRetry(func() (string, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(payload))
		if err != nil {
			return "", fmt.Errorf("creating openai request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("openai API call: %w", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("reading openai response: %w", err)
		}

		if resp.StatusCode >= 400 {
			return "", fmt.Errorf("openai API status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
		}

		var parsed chatCompletionResponse
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			return "", fmt.Errorf("parsing openai response: %w", err)
		}
		if parsed.Error != nil {
			return "", fmt.Errorf("openai API error: %s", parsed.Error.Message)
		}
		if len(parsed.Choices) == 0 {
			return "", fmt.Errorf("empty response from OpenAI")
		}

		return parsed.Choices[0].Message.Content, nil
	}, c.logger)
	if err != nil {
		return "", err
	}

	return extractJSON(result), nil
}

func extractJSON(s string) string {
	trimmed := strings.TrimSpace(s)

	if strings.HasPrefix(trimmed, "```") {
		trimmed = strings.TrimPrefix(trimmed, "```")
		trimmed = strings.TrimSpace(trimmed)
		if strings.HasPrefix(trimmed, "json") {
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "json"))
		}
		if idx := strings.LastIndex(trimmed, "```"); idx >= 0 {
			trimmed = strings.TrimSpace(trimmed[:idx])
		}
	}

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end >= start {
		return trimmed[start : end+1]
	}

	return trimmed
}

func withRetry[T any](fn func() (T, error), logger *slog.Logger) (T, error) {
	delays := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	var lastErr error

	for i := 0; i <= len(delays); i++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		lastErr = err

		if i < len(delays) {
			logger.Warn("retrying after error", "attempt", i+1, "error", err, "delay", delays[i])
			time.Sleep(delays[i])
		}
	}

	var zero T
	return zero, fmt.Errorf("all retries exhausted: %w", lastErr)
}
