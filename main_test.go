package traefik_responseheadersfilter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostResponsesHeaders(t *testing.T) {
	ctx := context.Background()
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		h := w.Header()

		h.Add("Header1", "value1")
		h.Add("Header2", "value2")
		h.Add("Header3", "value3")
		h.Add("Header4", "value4")

		_, _ = w.Write([]byte("foo"))
	})

	cfg := CreateConfig()
	cfg.Headers = []string{"header2", "header4"}
	mid, err := New(ctx, next, cfg, "testing-header-filter")
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(mid)
	t.Cleanup(server.Close)
	frontendClient := server.Client()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)

	res, err := frontendClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assertHeader(t, res, "Header1", "")
	assertHeader(t, res, "Header2", "value2")
	assertHeader(t, res, "Header3", "")
	assertHeader(t, res, "Header4", "value4")
}

func assertHeader(t *testing.T, res *http.Response, key, expected string) {
	t.Helper()

	if res.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", res.Header.Get(key))
	}
}
