package detect

import (
	"testing"

	"github.com/SteveSimpson/frontman/config"
	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupResponseStatusTestEnv(t *testing.T) (cleanup func()) {
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
		DetectResponseStatusExpireSeconds:   60,
		DetectResponseStatusStatusThreshold: 3,
		DetectResponseStatusIPThreshold:     5,
	}
	return func() {
		redisClient = oldRedis
		cfg = oldCfg
		rdb.Close()
		s.Close()
	}
}

func TestResponseStatusDetector_Detect_StatusThreshold(t *testing.T) {
	cleanup := setupResponseStatusTestEnv(t)
	defer cleanup()

	detector := &ResponseStatusDetector{}
	ip := "1.2.3.4"
	status := 404
	for i := 0; i < 3; i++ {
		detector.Detect(ip, status, nil)
	}
	score, err := detector.Detect(ip, status, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if score < 3 {
		t.Errorf("Expected score >= 3 for status threshold, got %d", score)
	}
	if len(detector.Details) == 0 {
		t.Error("Expected Details to be populated")
	}
}

func TestResponseStatusDetector_Detect_IPThreshold(t *testing.T) {
	cleanup := setupResponseStatusTestEnv(t)
	defer cleanup()

	detector := &ResponseStatusDetector{}
	ip := "5.6.7.8"
	status := 500
	for i := 0; i < 5; i++ {
		detector.Detect(ip, status, nil)
	}
	score, err := detector.Detect(ip, status, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if score < 5 {
		t.Errorf("Expected score >= 5 for IP threshold, got %d", score)
	}
	found := false
	for _, d := range detector.Details {
		if d.Match == "Total Errors" {
			found = true
		}
	}
	if !found {
		t.Error("Expected Details to include Total Errors for IP threshold")
	}
}

func TestResponseStatusDetector_Detect_BelowThreshold(t *testing.T) {
	cleanup := setupResponseStatusTestEnv(t)
	defer cleanup()

	detector := &ResponseStatusDetector{}
	ip := "9.9.9.9"
	status := 400
	score, err := detector.Detect(ip, status, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if score != 0 {
		t.Errorf("Expected score 0 below threshold, got %d", score)
	}
}

func TestResponseStatusDetector_GetDetails(t *testing.T) {
	detector := &ResponseStatusDetector{}
	detector.Details = DetectDetails{{Detector: "ResponseErrorDetector", Source: "1.2.3.4 Status Code: 404", Match: "Status Code 404", Value: 3}}
	details := detector.GetDetails()
	if len(details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(details))
	}
	if details[0].Detector != "ResponseErrorDetector" {
		t.Errorf("Expected Detector to be 'ResponseErrorDetector', got '%s'", details[0].Detector)
	}
}
