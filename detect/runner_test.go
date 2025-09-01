package detect

import (
	"net/http"
	"testing"

	"github.com/SteveSimpson/frontman/config"
	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupRunnerTestEnv(t *testing.T) (cleanup func()) {
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
		DetectBruteForcePath:             "/login",
		DetectBruteForceUsernameField:    "username",
		DetectBruteForcePasswordField:    "password",
		DetectBruteForceSalt:             "salt",
		DetectBruteForceExpireSeconds:    60,
		DetectBruteForceAlarmThreshold:   5,
		DetectSQLInjectionAlertThreshold: 7,
	}
	return func() {
		redisClient = oldRedis
		cfg = oldCfg
		rdb.Close()
		s.Close()
	}
}

func TestRequestDetectors_NoDetections(t *testing.T) {
	cleanup := setupRunnerTestEnv(t)
	defer cleanup()

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1"

	// Should not panic or error
	RequestDetectors(req)
}
