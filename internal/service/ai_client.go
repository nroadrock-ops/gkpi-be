package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type AIServiceClient interface {
	ForwardRequest(path string, payload interface{}) (map[string]interface{}, error)
}

type aiServiceClient struct {
	baseURL     string
	internalKey string
	client      *http.Client
}

func NewAIServiceClient(baseURL, internalKey string) AIServiceClient {
	return &aiServiceClient{
		baseURL:     baseURL,
		internalKey: internalKey,
		client: &http.Client{
			Timeout: 15 * time.Second, // Timeout gracefully
		},
	}
}

func (c *aiServiceClient) ForwardRequest(path string, payload interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader
	if payload != nil {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// WAJIB menyertakan header X-Internal-Key
	req.Header.Set("X-Internal-Key", c.internalKey)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("[AI Client Error] Failed to reach AI service at %s: %v", url, err)
		return nil, fmt.Errorf("AI service is unreachable or timed out")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("[AI Client Error] AI service returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("AI service returned error status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode AI response: %w", err)
	}

	return result, nil
}
