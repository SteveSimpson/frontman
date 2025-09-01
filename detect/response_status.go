package detect

import (
	"strconv"
	"time"
)

type ResponseStatusDetector struct {
	Details DetectDetails
}

func (d *ResponseStatusDetector) Detect(clientIP string, statusCode int, responseBody []byte) (int, error) {
	keyIp := "response_errors_ip"
	keyIpStatus := "response_errors_ip_status"

	returnCount := 0

	ipAndStatus := clientIP + ":" + strconv.Itoa(statusCode)

	if statusCode >= 400 {
		expiration := time.Duration(cfg.DetectResponseStatusExpireSeconds) * time.Second

		// Increment the count for this IP
		redisClient.HIncrBy(ctx, keyIp, clientIP, 1)
		redisClient.HExpire(ctx, keyIp, expiration, clientIP)

		// Increment the count for this IP and status code
		redisClient.HIncrBy(ctx, keyIpStatus, ipAndStatus, 1)
		redisClient.HExpire(ctx, keyIpStatus, expiration, ipAndStatus)

		ipStatusCount, err := redisClient.HGet(ctx, keyIpStatus, ipAndStatus).Int()
		if err != nil {
			return 0, err
		}

		// Check if the count for this IP and status code exceeds the threshold
		if ipStatusCount >= cfg.DetectResponseStatusStatusThreshold {
			returnCount = ipStatusCount
			d.Details = append(d.Details, struct {
				Detector string `json:"detector"`
				Source   string `json:"source"`
				Match    string `json:"match"`
				Value    int    `json:"value"`
			}{
				Detector: "ResponseErrorDetector",
				Source:   clientIP + " Status Code: " + strconv.Itoa(statusCode),
				Match:    "Status Code " + strconv.Itoa(statusCode),
				Value:    ipStatusCount,
			})
		}

		// Check if the count for this IP exceeds the threshold
		ipCount, err := redisClient.HGet(ctx, keyIp, clientIP).Int()
		if err != nil {
			return 0, err
		}
		if ipCount >= cfg.DetectResponseStatusIPThreshold {
			// ipCount should be equal to or greater than the ipStatusCount;
			// so return it if we are triggering on the IP threshold
			returnCount = ipCount
			d.Details = append(d.Details, struct {
				Detector string `json:"detector"`
				Source   string `json:"source"`
				Match    string `json:"match"`
				Value    int    `json:"value"`
			}{
				Detector: "ResponseErrorDetector",
				Source:   clientIP,
				Match:    "Total Errors",
				Value:    ipCount,
			})
		}
	}

	return returnCount, nil
}

func (d *ResponseStatusDetector) GetDetails() DetectDetails {
	return d.Details
}
