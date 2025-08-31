package detect

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type BruteForceDetector struct{}

var ctx = context.Background()

func (b *BruteForceDetector) Detect(req *http.Request) (int, error) {
	score := 0

	// only check if this is the login path
	if !strings.Contains(req.URL.Path, cfg.DetectBruteForcePath) {
		return 0, nil
	}

	score += b.incrementAndCheckRedis("bruteforce:ip", req.RemoteAddr)
	if req.FormValue(cfg.DetectBruteForceUsernameField) != "" {
		score += b.incrementAndCheckRedis("bruteforce:username", req.FormValue(cfg.DetectBruteForceUsernameField))
	}
	if req.FormValue(cfg.DetectBruteForcePasswordField) != "" {
		score += b.incrementAndCheckRedis("bruteforce:password", req.FormValue(cfg.DetectBruteForcePasswordField))
	}

	return score, nil
}

func (b *BruteForceDetector) incrementAndCheckRedis(key string, field string) int {
	// hash the field to avoid storing sensitive info in redis
	h := sha256.New()
	h.Write([]byte(field))
	hash := fmt.Sprintf("%x", h.Sum(nil))

	redisClient.HIncrBy(ctx, key, hash, 1)
	expiration := time.Duration(cfg.DetectBruteForceExpireSeconds) * time.Second

	redisClient.HExpire(ctx, key, expiration, hash)

	attempts, err := redisClient.HGet(ctx, key, hash).Int()
	if err != nil {
		// Not much we can do here, just log and return
		log.Printf("Error getting %s attempts from Redis: %v", key, err)
		return 0
	}

	if attempts >= cfg.DetectBruteForceAlarmThreshold {
		// Not really sure of the score here, just return the number of attempts for now
		return attempts
	}
	return 0
}
