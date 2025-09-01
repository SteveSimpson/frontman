package detect

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/SteveSimpson/frontman/config"
	"github.com/redis/go-redis/v9"
)

// DetectDetails holds details about detections made by detectors.
type DetectDetails []struct {
	Detector string `json:"detector"`
	Source   string `json:"source"`
	Match    string `json:"match"`
	Value    int    `json:"value"`
}

// DetectionReport holds a report of detections for a request.
type DetectionReport struct {
	ClientIP         string        `json:"client_ip"`
	Method           string        `json:"method"`
	RequestURI       string        `json:"request_uri"`
	UserAgent        string        `json:"user_agent"`
	Referrer         string        `json:"referrer"`
	TotalScore       int           `json:"total_score"`
	DetectorsTrigged int           `json:"detectors_triggered"`
	Details          DetectDetails `json:"details"`
}

type Detector interface {
	Detect(*http.Request) (int, error)
	GetDetails() DetectDetails
}

var redisClient *redis.Client
var cfg *config.ProxyConfig

func InitDetectors(c *config.ProxyConfig, r *redis.Client) {
	cfg = c
	redisClient = r
}

// RunDetectors runs all configured detectors on the given HTTP request.
func RequestDetectors(req *http.Request) {
	detectors := []Detector{
		&BruteForceDetector{},
		&SQLInjectionDetector{},
		&HostileSourceDetector{},
	}

	var combinedDetails DetectDetails

	err := req.ParseForm()
	if err != nil {
		// handle error
		log.Printf("Error parsing form: %v", err)
	}

	requestScore := 0
	detectorTriggres := 0

	for _, d := range detectors {
		detectorScore, err := d.Detect(req)

		if err != nil {
			// handle error
			log.Printf("Error running detector: %v", err)
			continue
		}
		if detectorScore > 0 {
			detectorTriggres++
			combinedDetails = append(combinedDetails, d.GetDetails()...)
		}
		requestScore += detectorScore
	}

	if requestScore > 0 {
		// Split IP from port in RemoteAddr
		clientIP := req.RemoteAddr
		if host, _, err := net.SplitHostPort(clientIP); err == nil {
			clientIP = host
		}

		report := DetectionReport{
			ClientIP:         clientIP,
			Method:           req.Method,
			RequestURI:       req.RequestURI,
			UserAgent:        req.UserAgent(),
			Referrer:         req.Referer(),
			TotalScore:       requestScore,
			DetectorsTrigged: detectorTriggres,
			Details:          combinedDetails,
		}

		jsonReport, err := json.Marshal(report)
		if err != nil {
			log.Printf("Error marshaling report: %v", err)
			return
		}

		// For now, just log the report. In a real system, you might send this to a logging service or alerting system.
		log.Printf("%s", jsonReport)
	}
}
