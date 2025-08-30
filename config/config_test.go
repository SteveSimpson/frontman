package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	os.Setenv("FRONTMAN_PROXY_ADDRESS", "127.0.0.1")
	os.Setenv("FRONTMAN_PROXY_PORT", "8080")
	os.Setenv("FRONTMAN_PRIVATE_URL", "http://private")
	os.Setenv("FRONTMAN_PUBLIC_URL", "http://public")
	defer func() {
		os.Unsetenv("FRONTMAN_PROXY_ADDRESS")
		os.Unsetenv("FRONTMAN_PROXY_PORT")
		os.Unsetenv("FRONTMAN_PRIVATE_URL")
		os.Unsetenv("FRONTMAN_PUBLIC_URL")
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.ProxyAddress != "127.0.0.1" {
		t.Errorf("expected ProxyAddress '127.0.0.1', got '%s'", cfg.ProxyAddress)
	}
	if cfg.ProxyPort != "8080" {
		t.Errorf("expected ProxyPort '8080', got '%s'", cfg.ProxyPort)
	}
	if cfg.PrivateURL != "http://private" {
		t.Errorf("expected PrivateURL 'http://private', got '%s'", cfg.PrivateURL)
	}
	if cfg.PublicURL != "http://public" {
		t.Errorf("expected PublicURL 'http://public', got '%s'", cfg.PublicURL)
	}
}

func TestLoadConfig_MissingEnv(t *testing.T) {
	os.Unsetenv("FRONTMAN_PROXY_ADDRESS")
	os.Unsetenv("FRONTMAN_PROXY_PORT")
	os.Unsetenv("FRONTMAN_PRIVATE_URL")
	os.Unsetenv("FRONTMAN_PUBLIC_URL")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error when env vars are missing, got nil")
	}
}
