package lcp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestFetchFrom_Slice(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `[{"name":"foo","value":1},{"name":"bar","value":2}]`)
		require.NoError(t, err)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, "dummy-token")
	result, err := FetchFrom[TestData](client, "/")
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, "foo", result[0].Name)
}

func TestFetchFrom_Map(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"a":{"name":"baz","value":42}}`)
		require.NoError(t, err)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, "dummy-token")
	result, err := FetchFrom[TestData](client, "/map")
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "baz", result[0].Name)
}

func TestFetchFrom_SingleObject(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"name":"solo","value":99}`)
		require.NoError(t, err)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, "dummy-token")
	result, err := FetchFrom[TestData](client, "/single")
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "solo", result[0].Name)
}

func TestFetchFrom_EnvelopeWithDataSlice(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"status":200,"message":"OK","data":[{"name":"env","value":5}]}`)
		require.NoError(t, err)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, "dummy-token")
	result, err := FetchFrom[TestData](client, "/envelope-slice")
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "env", result[0].Name)
}

func TestFetchFrom_EnvelopeWithDataMap(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"status":200,"message":"OK","data":{"x":{"name":"mapped","value":10}}}`)
		require.NoError(t, err)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, "dummy-token")
	result, err := FetchFrom[TestData](client, "/envelope-map")
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "mapped", result[0].Name)
}

func TestFetchFrom_EnvelopeWithErrorStatus(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"status":500,"message":"Internal error","data":null}`)
		require.NoError(t, err)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, "dummy-token")
	_, err := FetchFrom[TestData](client, "/error-envelope")
	require.Error(t, err)
}

func TestFetchFrom_InvalidJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"invalid_json":`)
		require.NoError(t, err)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, "dummy-token")
	_, err := FetchFrom[TestData](client, "/invalid-json")
	require.Error(t, err)
}

func TestFetchFrom_APIErrorStatus(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, err := io.WriteString(w, "forbidden\n")
		require.NoError(t, err)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, "dummy-token")
	_, err := FetchFrom[TestData](client, "/forbidden")
	require.Error(t, err)
}

func TestFetchOneFrom(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, `{"name":"one","value":7}`)
		require.NoError(t, err)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(server.URL, "dummy-token")
	result, err := FetchOneFrom[TestData](client, "/one")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "one", result.Name)
}
