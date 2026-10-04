package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Provider defines the minimal interface we need from an LLM.
type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// OpenAIProvider wraps the OpenAI REST API.
type OpenAIProvider struct {
	APIKey    string
	Model     string
	HTTPCli   *http.Client
	Endpoint  string // e.g., https://api.openai.com/v1/chat/completions
}

// NewOpenAIProvider creates a provider using the API key from viper.
func NewOpenAIProvider(v *config.Viper) (*OpenAIProvider, error) {
	key := v.GetString("llm_router.providers.0.api_key_env")
	if key == "" {
		// In a real build you would read the actual secret file or env var.
		// For the skeleton we accept an empty key and will return an error on use.
	}
	return &OpenAIProvider{
		APIKey:    key,
		Model:     v.GetString("llm_router.providers.0.model"),
		HTTPCli:   &http.Client{Timeout: 30 * time.Second},
		Endpoint:  "https://api.openai.com/v1/chat/completions",
	}, nil
}

// Generate sends a chat‑completion request to OpenAI.
func (p *OpenAIProvider) Generate(ctx context.Context, prompt string) (string, error) {
	if p.APIKey == "" {
		return "", errors.New("openai API key not configured")
	}
	payload := map[string]interface{}{
		"model": p.Model,
		"messages": []map[string]string{
			{ "role": "system", "content": "You are a helpful security assistant." },
			{ "role": "user",   "content": prompt },
		},
		"temperature": 0.2,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", p.Endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := p.HTTPCli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai error: %d", resp.StatusCode)
	}
	var respData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}
	if len(respData.Choices) == 0 {
		return "", errors.New("empty response from OpenAI")
	}
	return respData.Choices[0].Message.Content, nil
}

// OllamaProvider talks to a locally‑running Ollama instance.
type OllamaProvider struct {
	HTTPCli   *http.Client
	Endpoint  string // e.g., http://ollama:11434
	Model     string
}

func NewOllamaProvider(v *config.Viper) (*OllamaProvider, error) {
	return &OllamaProvider{
		HTTPCli: &http.Client{Timeout: 60 * time.Second},
		Endpoint: v.GetString("llm_router.providers.1.endpoint"),
		Model:    v.GetString("llm_router.providers.1.model"),
	}, nil
}

func (p *OllamaProvider) Generate(ctx context.Context, prompt string) (string, error) {
	payload := map[string]interface{}{
		"model":  p.Model,
		"prompt": prompt,
		"stream": false,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", p.Endpoint+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.HTTPCli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama error: %d", resp.StatusCode)
	}
	var respData struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}
	return respData.Response, nil
}
