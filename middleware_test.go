package traefikgeoip2_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	traefikgeoip2 "github.com/7i11/traefik-geoip2"
)

func TestTraefikGeoip2(t *testing.T) {
	cfg := traefikgeoip2.CreateConfig()
	cfg.FromHeader = "X-Forwarded-For"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefikgeoip2.New(ctx, next, cfg, "traefikgeoip2")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, "http://localhost", nil,
	)
	req.Header[cfg.FromHeader] = []string{"8.8.8.8"}

	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Geoip2-Country", "US")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
