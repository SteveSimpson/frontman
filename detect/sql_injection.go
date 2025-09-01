package detect

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// This detector looks for SQL injection attempts in request parameters.
type SQLInjectionDetector struct {
	Details DetectDetails
}

// Detect looks for SQL injection attempts in request parameters.
func (s *SQLInjectionDetector) Detect(req *http.Request) (int, error) {
	score := 0

	patterns := map[string]int{
		"UNION SELECT":       5,
		"DROP TABLE":         5,
		"DELETE FROM":        5,
		"' OR '1'='1":        7,
		"\" OR \"1\"=\"1":    7,
		"--":                 3,
		";":                  2,
		"/*":                 2,
		"*/":                 2,
		"xp_cmdshell":        5,
		"exec":               4,
		"INFORMATION_SCHEMA": 4,
		"sysobjects":         4,
		"@@version":          3,
		"LOAD_FILE":          5,
		"INTO OUTFILE":       5,
		"CONCAT":             2,
		"GROUP BY":           2,
		"HAVING":             2,
		"SLEEP":              4,
		"BENCHMARK":          4,
		"WAITFOR DELAY":      4,
	}

	// Check query parameters
	queryParams := req.URL.Query()
	for _, values := range queryParams {
		for _, value := range values {
			for pattern, points := range patterns {
				value, _ = url.QueryUnescape(value)
				if containsIgnoreCase(value, pattern) {
					score += points
					s.Details = append(s.Details, struct {
						Detector string `json:"detector"`
						Source   string `json:"source"`
						Match    string `json:"match"`
						Value    int    `json:"value"`
					}{
						Detector: "SQLInjectionDetector",
						Source:   "query parameter",
						Match:    pattern,
						Value:    points,
					})
				}
			}
		}
	}

	// Check form data (skipped for GET requests to avoid double-checking query params)
	if req.Method != "GET" {
		if err := req.ParseForm(); err == nil {
			for _, values := range req.Form {
				for _, value := range values {
					for pattern, points := range patterns {
						if containsIgnoreCase(value, pattern) {
							score += points
							s.Details = append(s.Details, struct {
								Detector string `json:"detector"`
								Source   string `json:"source"`
								Match    string `json:"match"`
								Value    int    `json:"value"`
							}{
								Detector: "SQLInjectionDetector",
								Source:   "form data",
								Match:    pattern,
								Value:    points,
							})
						}
					}
				}
			}
		} else {
			// If form parsing fails, we skip form data checks.
			fmt.Println("Failed to parse form data: ", err)
		}
	}

	// Don't alert unless we hit the threshold
	if score < cfg.DetectSQLInjectionAlertThreshold {
		return 0, nil
	}

	return score, nil
}

// containsIgnoreCase checks if substr is in s, ignoring case and spaces.
func containsIgnoreCase(s, substr string) bool {
	// normalize by remvoving spaces and making upper case
	s = strings.ToUpper(strings.ReplaceAll(s, " ", ""))
	substr = strings.ToUpper(strings.ReplaceAll(substr, " ", ""))

	return strings.Contains(s, substr)
}

func (s *SQLInjectionDetector) GetDetails() DetectDetails {
	return s.Details
}
