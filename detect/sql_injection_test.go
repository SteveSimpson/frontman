package detect

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/SteveSimpson/frontman/config"
)

func setupSQLTestEnv() (cleanup func()) {
	oldCfg := cfg
	cfg = &config.ProxyConfig{
		DetectSQLInjectionAlertThreshold: 7,
	}
	return func() { cfg = oldCfg }
}

func TestSQLInjectionDetector_Detect_BelowThreshold(t *testing.T) {
	cleanup := setupSQLTestEnv()
	defer cleanup()

	detector := &SQLInjectionDetector{}
	q := url.Values{}
	q.Set("foo", ";--") // score = 5 (2+3), below threshold 7
	req, _ := http.NewRequest("GET", "/?"+q.Encode(), nil)

	score, err := detector.Detect(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if score != 0 {
		t.Errorf("Expected score 0 below threshold, got %d", score)
	}
}

func TestSQLInjectionDetector_Detect_QueryParam(t *testing.T) {
	cleanup := setupSQLTestEnv()
	defer cleanup()

	detector := &SQLInjectionDetector{}
	q := url.Values{}
	q.Set("search", "foo; UNION SELECT bar")
	req, _ := http.NewRequest("GET", "/?"+q.Encode(), nil)

	score, err := detector.Detect(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if score == 0 {
		t.Errorf("Expected non-zero score for SQLi pattern, got %d", score)
	}
	if len(detector.Details) == 0 {
		t.Error("Expected Details to be populated")
	}
}

func TestSQLInjectionDetector_Detect_FormData(t *testing.T) {
	cleanup := setupSQLTestEnv()
	defer cleanup()

	detector := &SQLInjectionDetector{}
	form := url.Values{}
	form.Set("username", "admin")
	form.Set("password", "' OR '1'='1")
	reqBody := strings.NewReader(form.Encode())
	req, _ := http.NewRequest("POST", "/login", reqBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	score, err := detector.Detect(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if score == 0 {
		t.Errorf("Expected non-zero score for SQLi pattern, got %d", score)
	}
	if len(detector.Details) == 0 {
		t.Error("Expected Details to be populated")
	}
}

func TestSQLInjectionDetector_GetDetails(t *testing.T) {
	detector := &SQLInjectionDetector{}
	detector.Details = DetectDetails{{Detector: "SQLInjectionDetector", Source: "query parameter", Match: "UNION SELECT", Value: 5}}
	details := detector.GetDetails()
	if len(details) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(details))
	}
	if details[0].Detector != "SQLInjectionDetector" {
		t.Errorf("Expected Detector to be 'SQLInjectionDetector', got '%s'", details[0].Detector)
	}
}
