package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	os.Setenv("FRONTMAN_DETECT_BRUTEFORCE_PATH", "/login")
	os.Setenv("FRONTMAN_PROXY_ADDRESS", "127.0.0.1")
	os.Setenv("FRONTMAN_PROXY_PORT", "8080")
	os.Setenv("FRONTMAN_PRIVATE_URL", "http://private")
	os.Setenv("FRONTMAN_PUBLIC_URL", "http://public")
	os.Setenv("FRONTMAN_REDIS_HOST", "redis:6379")
	os.Setenv("FRONTMAN_DETECT_BRUTEFORCE_SALT", "random_salt_value")
	defer func() {
		os.Unsetenv("FRONTMAN_DETECT_BRUTEFORCE_PATH")
		os.Unsetenv("FRONTMAN_PROXY_ADDRESS")
		os.Unsetenv("FRONTMAN_PROXY_PORT")
		os.Unsetenv("FRONTMAN_PRIVATE_URL")
		os.Unsetenv("FRONTMAN_PUBLIC_URL")
		os.Unsetenv("FRONTMAN_REDIS_HOST")
		os.Unsetenv("FRONTMAN_DETECT_BRUTEFORCE_SALT")
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
	if cfg.DetectBruteForcePath != "/login" {
		t.Errorf("expected DetectBruteForcePath '/login', got '%s'", cfg.DetectBruteForcePath)
	}
	if cfg.DetectBruteForceUsernameField != "username" {
		t.Errorf("expected DetectBruteForceUsernameField 'username', got '%s'", cfg.DetectBruteForceUsernameField)
	}
	if cfg.DetectBruteForcePasswordField != "password" {
		t.Errorf("expected DetectBruteForcePasswordField 'password', got '%s'", cfg.DetectBruteForcePasswordField)
	}
	if cfg.DetectBruteForceAlarmThreshold != 10 {
		t.Errorf("expected DetectBruteForceAlarmThreshold 10, got %d", cfg.DetectBruteForceAlarmThreshold)
	}
	if cfg.DetectBruteForceExpireSeconds != 3600 {
		t.Errorf("expected DetectBruteForceExpireSeconds 3600, got %d", cfg.DetectBruteForceExpireSeconds)
	}
	if cfg.RedisHost != "redis:6379" {
		t.Errorf("expected RedisHost 'redis:6379', got '%s'", cfg.RedisHost)
	}
	if cfg.RedisPassword != "" {
		t.Errorf("expected RedisPassword '', got '%s'", cfg.RedisPassword)
	}
	if cfg.RedisDB != 0 {
		t.Errorf("expected RedisDB 0, got %d", cfg.RedisDB)
	}
	if cfg.DetectBruteForceSalt != "random_salt_value" {
		t.Error("expected DetectBruteForceSalt to be set, got empty string")
	}
	if cfg.DetectSQLInjectionAlertThreshold != 7 {
		t.Errorf("expected DetectSQLInjectionAlertThreshold 5, got %d", cfg.DetectSQLInjectionAlertThreshold)
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
