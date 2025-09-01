package detect

import (
	"context"
	"net/http"
	"testing"

	"github.com/SteveSimpson/frontman/config"
	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupHostileTestEnv(t *testing.T) (cleanup func()) {
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
	cfg = &config.ProxyConfig{}
	return func() {
		redisClient = oldRedis
		cfg = oldCfg
		rdb.Close()
		s.Close()
	}
}

func TestHostileSourceDetector_Detect_Hostile(t *testing.T) {
	cleanup := setupHostileTestEnv(t)
	defer cleanup()

	detector := &HostileSourceDetector{}
	ip := "1.2.3.4"
	// Add IP to hostile_ips set in miniredis
	err := redisClient.SAdd(context.TODO(), "hostile_ips", ip).Err()
	if err != nil {
		t.Fatalf("failed to add IP to hostile_ips: %v", err)
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = ip

	score, err := detector.Detect(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if score != 10 {
		t.Errorf("Expected score 10 for hostile IP, got %d", score)
	}
	if len(detector.Details) == 0 {
		t.Error("Expected Details to be populated")
	}
}

func TestHostileSourceDetector_Detect_NotHostile(t *testing.T) {
	cleanup := setupHostileTestEnv(t)
	defer cleanup()

	detector := &HostileSourceDetector{}
	ip := "5.6.7.8"

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = ip

	score, err := detector.Detect(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if score != 0 {
		t.Errorf("Expected score 0 for non-hostile IP, got %d", score)
	}
	if len(detector.Details) != 0 {
		t.Error("Expected Details to be empty for non-hostile IP")
	}
}

func TestHostileSourceDetector_GetDetails(t *testing.T) {
	detector := &HostileSourceDetector{}
	detector.Details = DetectDetails{{Detector: "HostileSourceDetector", Source: "IP Address", Match: "1.2.3.4", Value: 10}}
	details := detector.GetDetails()
	if len(details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(details))
	}
	if details[0].Detector != "HostileSourceDetector" {
		t.Errorf("Expected Detector to be 'HostileSourceDetector', got '%s'", details[0].Detector)
	}
}
