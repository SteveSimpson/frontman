package detect

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// BruteForceDetector looks for brute force login attempts based on IP, username, and password.
type BruteForceDetector struct {
	Details DetectDetails
}

var ctx = context.Background()

// Detect looks for brute force login attempts based on IP, username, and password.
func (b *BruteForceDetector) Detect(req *http.Request) (int, error) {
	score := 0

	// only check if this is the login path
	if !strings.Contains(req.URL.Path, cfg.DetectBruteForcePath) {
		return 0, nil
	}

	// Split IP from port in RemoteAddr
	clientIP := req.RemoteAddr
	if host, _, err := net.SplitHostPort(clientIP); err == nil {
		clientIP = host
	}

	score += b.incrementAndCheckRedis("bruteforce:ip", clientIP, clientIP)
	if req.FormValue(cfg.DetectBruteForceUsernameField) != "" {
		username := req.FormValue(cfg.DetectBruteForceUsernameField)
		score += b.incrementAndCheckRedis("bruteforce:username", username, username)
	}
	if req.FormValue(cfg.DetectBruteForcePasswordField) != "" {
		password := req.FormValue(cfg.DetectBruteForcePasswordField)
		score += b.incrementAndCheckRedis("bruteforce:password", password, sanitizePassword(password))
	}

	return score, nil
}

// incrementAndCheckRedis increments the count for the given key and field in Redis,
// sets the expiration, and checks if the count exceeds the alarm threshold.
// Returns a score if the threshold is exceeded, otherwise returns 0.
// details is the string to log (ip, username, or sanitized password)
// Logs are only made when the threshold is exceeded, so I am less concerned about logging sensitive info.
func (b *BruteForceDetector) incrementAndCheckRedis(key, field, details string) int {
	// hash the field to avoid storing sensitive info in redis
	h := sha256.New()
	h.Write([]byte(cfg.DetectBruteForceSalt + field))
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
		b.Details = append(b.Details, struct {
			Detector string `json:"detector"`
			Source   string `json:"source"`
			Match    string `json:"match"`
			Value    int    `json:"value"`
		}{
			Detector: "BruteForceDetector",
			Source:   key,
			Match:    details,
			Value:    attempts,
		})
		// Not really sure of the score here, just return the number of attempts for now
		return attempts
	}
	return 0
}

// sanitizePassword replaces all but the first and last character of a password with asterisks
func sanitizePassword(p string) string {
	if len(p) <= 2 {
		return p
	}
	return p[:1] + strings.Repeat("*", len(p)-2) + p[len(p)-1:]
}

func (b *BruteForceDetector) GetDetails() DetectDetails {
	return b.Details
}
