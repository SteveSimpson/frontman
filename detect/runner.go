package detect

import (
	"log"
	"net/http"

	"github.com/SteveSimpson/frontman/config"
	"github.com/redis/go-redis/v9"
)

type Detector interface {
	Detect(*http.Request) (int, error)
}

var redisClient *redis.Client
var cfg *config.ProxyConfig

func InitDetectors(c *config.ProxyConfig, r *redis.Client) {
	cfg = c
	redisClient = r
}

func RunDetectors(req *http.Request) {
	detectors := []Detector{
		&BruteForceDetector{},
		// &SQLInjectionDetector{},
	}

	err := req.ParseForm()
	if err != nil {
		// handle error
		log.Printf("Error parsing form: %v", err)
	}

	requestScore := 0
	detectorTriggres := 0

	log.Println(req.Form)

	for _, d := range detectors {
		detectorScore, err := d.Detect(req)

		if err != nil {
			// handle error
			log.Printf("Error running detector: %v", err)
			continue
		}
		if detectorScore > 0 {
			detectorTriggres++
		}
		requestScore += detectorScore
	}

	if requestScore > 0 {
		// log potential attack
		// probabluy want to do something more here (json to log collector?)
		log.Printf("%s for %s from %s scored %d with %d detectors triggered", req.Method, req.RequestURI, req.Host, requestScore, detectorTriggres)
	}
}
