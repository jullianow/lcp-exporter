package lcp

import (
	"fmt"
	"io"
	"net/http"
	"time"
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

func (c *Client) MakeRequest(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.BaseURL, endpoint), nil)
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
		return resp, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
