package lcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jullianow/lcp-exporter/internal"
)

type Client struct {
	BaseURL     string
	BearerToken string
	Client      *http.Client
}

func NewClient(baseURL, bearerToken string) *Client {
	return &Client{
		BaseURL:     baseURL,
		BearerToken: bearerToken,
		Client:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) MakeRequest(path string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		if cerr := resp.Body.Close(); cerr != nil {
			internal.LogWarn("MakeRequest", "Error closing response body: %v", cerr)
		}
		internal.LogWarn("MakeRequest", "Non-2xx response from %s: %d - %s", url, resp.StatusCode, string(body))
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func ParseEnvelope[T any](body []byte) ([]T, error) {
	var envelope struct {
		Status  int             `json:"status"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Data != nil {
		if envelope.Status != 0 && envelope.Status != http.StatusOK {
			internal.LogWarn("ParseEnvelope", "API error status: %d - %s", envelope.Status, envelope.Message)
			return nil, fmt.Errorf("API error %d: %s", envelope.Status, envelope.Message)
		}

		var dataSlice []T
		if err := json.Unmarshal(envelope.Data, &dataSlice); err == nil {
			return dataSlice, nil
		}

		var dataMap map[string]T
		if err := json.Unmarshal(envelope.Data, &dataMap); err == nil {
			var slice []T
			for _, v := range dataMap {
				slice = append(slice, v)
			}
			return slice, nil
		}

		var single T
		if err := json.Unmarshal(envelope.Data, &single); err == nil {
			return []T{single}, nil
		}
	}

	var result []T
	if err := json.Unmarshal(body, &result); err == nil {
		return result, nil
	}

	var resultMap map[string]T
	if err := json.Unmarshal(body, &resultMap); err == nil {
		var slice []T
		for _, v := range resultMap {
			slice = append(slice, v)
		}
		return slice, nil
	}

	var single T
	if err := json.Unmarshal(body, &single); err == nil {
		return []T{single}, nil
	}

	internal.LogError("ParseEnvelope", "Failed to unmarshal response body: %s", string(body))
	return nil, fmt.Errorf("failed to parse response")
}

func FetchFrom[T any](c *Client, path string) ([]T, error) {
	resp, err := c.MakeRequest(path)
	if err != nil {
		internal.LogError("FetchFrom", "Request failed for path %s: %v", path, err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			internal.LogWarn("FetchFrom", "Error closing response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		internal.LogError("FetchFrom", "Failed to read body from path %s: %v", path, err)
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	return ParseEnvelope[T](body)
}

func FetchOneFrom[T any](c *Client, path string) (*T, error) {
	results, err := FetchFrom[T](c, path)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}
	return &results[0], nil
}
