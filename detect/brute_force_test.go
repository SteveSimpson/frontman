package detect

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/SteveSimpson/frontman/config"
	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// Setup miniredis and go-redis for integration-style test
func setupTestEnv(t *testing.T) (cleanup func()) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
		DB:   0,
	})
	oldRedis := redisClient
	oldCfg := cfg
	redisClient = rdb
	cfg = &config.ProxyConfig{
		DetectBruteForcePath:           "/login",
		DetectBruteForceUsernameField:  "username",
		DetectBruteForcePasswordField:  "password",
		DetectBruteForceSalt:           "salt",
		DetectBruteForceExpireSeconds:  60,
		DetectBruteForceAlarmThreshold: 5,
	}
	return func() {
		redisClient = oldRedis
		cfg = oldCfg
		rdb.Close()
		s.Close()
	}
}

func TestBruteForceDetector_Detect(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	detector := &BruteForceDetector{}
	form := url.Values{}
	form.Set("username", "user1")
	form.Set("password", "pass1")
	req, _ := http.NewRequest("POST", "/login", nil)
	req.URL.Path = "/login"
	req.RemoteAddr = "1.2.3.4:12345"
	req.Form = form

	// Call Detect 5 times to exceed the threshold (5)
	for i := 0; i < 5; i++ {
		detector.Detect(req)
	}
	score, err := detector.Detect(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if score < 5 {
		t.Errorf("Expected score >= 5, got %d", score)
	}
	if len(detector.Details) == 0 {
		t.Error("Expected Details to be populated")
	}
}

func TestBruteForceDetector_GetDetails(t *testing.T) {
	detector := &BruteForceDetector{}
	detector.Details = DetectDetails{{Detector: "BruteForceDetector", Source: "ip", Match: "1.2.3.4", Value: 5}}
	details := detector.GetDetails()
	if len(details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(details))
	}
	if details[0].Detector != "BruteForceDetector" {
		t.Errorf("Expected Detector to be 'BruteForceDetector', got '%s'", details[0].Detector)
	}
}
