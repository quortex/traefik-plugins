// Package responseheadersfilter contains the implementation of a middleware that filters unwanted response headers.
package traefik_responseheadersfilter

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	Headers []string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: []string{},
	}
}

// responseheadersfilter is a middleware that filters unwanted response headers.
type responseheadersfilter struct {
	next    http.Handler
	name    string
	headers []string
}

// New creates a new instance of the responseheadersfilter middleware.
// It takes a context.Context, an http.Handler, a *Config, and a name string as parameters.
// It returns an http.Handler and an error.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	infoLogger := log.New(io.Discard, "INFO: responseheadersfilter: ", log.Ldate|log.Ltime)
	infoLogger.SetOutput(os.Stdout)
	infoLogger.Printf("Headers allowed: %s", config.Headers)

	return &responseheadersfilter{
		headers: config.Headers,
		next:    next,
		name:    name,
	}, nil
}

// ServeHTTP is the method that handles the HTTP request.
// It calls the next handler in the chain and modifies the response headers.
func (f *responseheadersfilter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if f.next != nil {
		f.next.ServeHTTP(newResponseModifier(rw, req, f.PostRequestDeleteResponseHeaders), req)
	}
}

// PostRequestDeleteResponseHeaders is a method that deletes unwanted response headers.
// It is called AFTER the response is generated from the backend.
func (f *responseheadersfilter) PostRequestDeleteResponseHeaders(res *http.Response) error {
	for key := range res.Header {
		found := false
		for _, header := range f.headers {
			found = found || strings.EqualFold(key, header)
		}
		if !found {
			res.Header.Del(key)
		}
	}
	return nil
}
