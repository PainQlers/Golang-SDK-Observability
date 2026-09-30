package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrapHTTPClientPerformsRequest(t *testing.T) {
	var receivedPath string
	var sawTraceparent bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		sawTraceparent = r.Header.Get("traceparent") != ""
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := WrapHTTPClient(server.Client())
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/hello", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
	}
	if receivedPath != "/hello" {
		t.Fatalf("expected path /hello, got %s", receivedPath)
	}
	if !sawTraceparent {
		t.Fatal("expected trace context to be injected into the outbound request")
	}
}
