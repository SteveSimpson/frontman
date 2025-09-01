package detect

// This detector simply looks for requests from known hostile sources.
// Since I already used redis, I am using that here.
// I would probably use a database where I could look at up networks and ranges.
// I would also want to use a proper geoip database to look for suspicious locations,
// but that could be done on the database side as well.

import (
	"net"
	"net/http"
)

// HostileSourceDetector checks if the request comes from a known hostile source.
type HostileSourceDetector struct {
	Details DetectDetails
}

// Detect checks if the request comes from a known hostile source.
func (d *HostileSourceDetector) Detect(req *http.Request) (int, error) {

	// Split IP from port in RemoteAddr
	clientIP := req.RemoteAddr
	if host, _, err := net.SplitHostPort(clientIP); err == nil {
		clientIP = host
	}

	isHostile, err := redisClient.SIsMember(req.Context(), "hostile_ips", clientIP).Result()
	if err != nil {
		return 0, err
	}
	if isHostile {
		d.Details = append(d.Details, struct {
			Detector string `json:"detector"`
			Source   string `json:"source"`
			Match    string `json:"match"`
			Value    int    `json:"value"`
		}{
			Detector: "HostileSourceDetector",
			Source:   "IP Address",
			Match:    clientIP,
			Value:    10,
		})

		return 10, nil // High score for known hostile source
	}
	return 0, nil
}

func (d *HostileSourceDetector) GetDetails() DetectDetails {
	return d.Details
}
