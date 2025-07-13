package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"study_buddy/internal/config"
)

type Client interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type OpenRouterClient struct {
	APIKey string
	OpenAi *config.OpenAI
}

func NewOpenRouterClient(openAi *config.OpenAI) *OpenRouterClient {
	return &OpenRouterClient{
		APIKey: openAi.ApiKeys[openAi.Ind],
		OpenAi: openAi,
	}
}

func (c *OpenRouterClient) Complete(ctx context.Context, prompt string) (string, error) {
	body := map[string]interface{}{
		"model": "deepseek/deepseek-r1-0528:free",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	bodyJSON, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(bodyJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm error: %s", string(bodyBytes))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("llm response had no choices")
	}

	c.OpenAi.Ind = (c.OpenAi.Ind + 1) % int32(len(c.OpenAi.ApiKeys))
	c.APIKey = c.OpenAi.ApiKeys[c.OpenAi.Ind]

	return result.Choices[0].Message.Content, nil
}
