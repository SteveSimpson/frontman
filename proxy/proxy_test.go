package proxy

import (
	"net/http/httptest"
	"testing"

	"github.com/SteveSimpson/frontman/config"
)

func TestNewProxy(t *testing.T) {
	cfg := &config.ProxyConfig{
		PublicURL:  "http://example.com",
		PrivateURL: "http://localhost:8080",
	}
	p, err := NewProxy(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if p.publicUrl.String() != cfg.PublicURL {
		t.Errorf("Expected publicUrl %s, got %s", cfg.PublicURL, p.publicUrl.String())
	}
	if p.privateUrl.String() != cfg.PrivateURL {
		t.Errorf("Expected privateUrl %s, got %s", cfg.PrivateURL, p.privateUrl.String())
	}
	if p.proxy == nil {
		t.Error("Expected proxy to be initialized")
	}
}

func TestHandleRequest(t *testing.T) {
	cfg := &config.ProxyConfig{
		PublicURL:  "http://example.com",
		PrivateURL: "http://localhost:8080",
	}
	p, err := NewProxy(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// This will call the reverse proxy, which will try to reach the private URL.
	// Since there's no server running, we just want to make sure it doesn't panic.
	// The error will be written to the ResponseWriter.
	p.HandleRequest(rec, req)

	// Check that the request was rewritten
	if req.URL.Host != "example.com" {
		t.Errorf("Expected URL.Host to be 'example.com', got '%s'", req.URL.Host)
	}
	if req.URL.Scheme != "http" {
		t.Errorf("Expected URL.Scheme to be 'http', got '%s'", req.URL.Scheme)
	}
	if req.Host != "example.com" {
		t.Errorf("Expected Host to be 'example.com', got '%s'", req.Host)
	}
}
