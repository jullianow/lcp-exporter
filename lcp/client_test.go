package lcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-endpoint", r.URL.Path)

		_, err := w.Write([]byte(`{"message": "success"}`))
		if err != nil {
			t.Fatalf("Error writing response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "")

	resp, err := client.MakeRequest("/test-endpoint")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.Equal(t, "success", responseData["message"])
}
