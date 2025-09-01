package config

import (
	"errors"
	"os"
	"strconv"
)

type ProxyConfig struct {
	DetectBruteForcePath                string `json:"detect_brute_force_path"`
	DetectBruteForceUsernameField       string `json:"detect_brute_force_username_field"`
	DetectBruteForcePasswordField       string `json:"detect_brute_force_password_field"`
	DetectBruteForceAlarmThreshold      int    `json:"detect_brute_force_alarm_threshold"`
	DetectBruteForceExpireSeconds       int    `json:"detect_brute_force_expire_seconds"`
	DetectBruteForceSalt                string `json:"detect_brute_force_salt"`
	DetectResponseStatusExpireSeconds   int    `json:"detect_response_status_expire_seconds"`
	DetectResponseStatusIPThreshold     int    `json:"detect_response_status_ip_threshold"`
	DetectResponseStatusStatusThreshold int    `json:"detect_response_status_status_threshold"`
	DetectSQLInjectionAlertThreshold    int    `json:"detect_sql_injection_alert_threshold"`
	ProxyAddress                        string `json:"proxy_address"`
	ProxyPort                           string `json:"proxy_port"`
	PrivateURL                          string `json:"private_url"`
	PublicURL                           string `json:"public_url"`
	RedisHost                           string `json:"redis_host"`
	RedisPassword                       string `json:"redis_password"`
	RedisDB                             int    `json:"redis_db"`
}

func LoadConfig() (*ProxyConfig, error) {

	proxyAddress, found := os.LookupEnv("FRONTMAN_PROXY_ADDRESS")
	if !found {
		return nil, errors.New("FRONTMAN_PROXY_ADDRESS not set in environment")
	}

	proxyPort, found := os.LookupEnv("FRONTMAN_PROXY_PORT")
	if !found {
		return nil, errors.New("FRONTMAN_PROXY_PORT not set in environment")
	}

	privateURL, found := os.LookupEnv("FRONTMAN_PRIVATE_URL")
	if !found {
		return nil, errors.New("FRONTMAN_PRIVATE_URL not set in environment")
	}

	publicURL, found := os.LookupEnv("FRONTMAN_PUBLIC_URL")
	if !found {
		return nil, errors.New("FRONTMAN_PUBLIC_URL not set in environment")
	}

	bruteForcePath := os.Getenv("FRONTMAN_DETECT_BRUTEFORCE_PATH")

	bruteForceUsernameField, found := os.LookupEnv("FRONTMAN_DETECT_BRUTEFORCE_USERNAME_FIELD")
	if !found {
		bruteForceUsernameField = "username"
	}
	bruteForcePasswordField, found := os.LookupEnv("FRONTMAN_DETECT_BRUTEFORCE_PASSWORD_FIELD")
	if !found {
		bruteForcePasswordField = "password"
	}

	bruteForceAlarmThresholdStr, found := os.LookupEnv("FRONTMAN_DETECT_BRUTEFORCE_ALARM_THRESHOLD")
	if !found {
		bruteForceAlarmThresholdStr = "10"
	}

	bruteForceAlarmThreshold, err := strconv.Atoi(bruteForceAlarmThresholdStr)
	if err != nil {
		return nil, errors.New("FRONTMAN_DETECT_BRUTEFORCE_ALARM_THRESHOLD must be an integer")
	}

	bruteForceExpireSecondsStr, found := os.LookupEnv("FRONTMAN_DETECT_BRUTEFORCE_EXPIRE_SECONDS")
	if !found {
		bruteForceExpireSecondsStr = "3600" // default to 60 minutes
	}

	bruteForceExpireSeconds, err := strconv.Atoi(bruteForceExpireSecondsStr)
	if err != nil {
		return nil, errors.New("FRONTMAN_DETECT_BRUTEFORCE_EXPIRE_SECONDS must be an integer")
	}

	bruteForceSalt, found := os.LookupEnv("FRONTMAN_DETECT_BRUTEFORCE_SALT")
	if !found {
		return nil, errors.New("FRONTMAN_DETECT_BRUTEFORCE_SALT not set in environment")
	}

	redisHost, found := os.LookupEnv("FRONTMAN_REDIS_HOST")
	if !found {
		return nil, errors.New("FRONTMAN_REDIS_HOST not set in environment")
	}

	redisPassword, found := os.LookupEnv("FRONTMAN_REDIS_PASSWORD")
	if !found {
		redisPassword = "" // No password by default for local development
	}

	redisDB := 0 // Default DB
	redisDBStr, found := os.LookupEnv("FRONTMAN_REDIS_DB")
	if found {
		// Convert to int
		var err error
		redisDB, err = strconv.Atoi(redisDBStr)
		if err != nil {
			return nil, errors.New("FRONTMAN_REDIS_DB must be an integer")
		}
	}

	responseStatusExpireSecondsStr, found := os.LookupEnv("FRONTMAN_DETECT_RESPONSE_STATUS_EXPIRE_SECONDS")
	if !found {
		responseStatusExpireSecondsStr = "300" // default to 5 minutes
	}
	responseStatusExpireSeconds, err := strconv.Atoi(responseStatusExpireSecondsStr)
	if err != nil {
		return nil, errors.New("FRONTMAN_DETECT_RESPONSE_STATUS_EXPIRE_SECONDS must be an integer")
	}

	responseStatusIPThresholdStr, found := os.LookupEnv("FRONTMAN_DETECT_RESPONSE_STATUS_IP_THRESHOLD")
	if !found {
		responseStatusIPThresholdStr = "50" // default to 50 errors (for all errors for an IP)
	}
	responseStatusIPThreshold, err := strconv.Atoi(responseStatusIPThresholdStr)
	if err != nil {
		return nil, errors.New("FRONTMAN_DETECT_RESPONSE_STATUS_IP_THRESHOLD must be an integer")
	}

	responseStatusStatusThresholdStr, found := os.LookupEnv("FRONTMAN_DETECT_RESPONSE_STATUS_STATUS_THRESHOLD")
	if !found {
		responseStatusStatusThresholdStr = "10" // default to 10 errors per status code per IP
	}
	responseStatusStatusThreshold, err := strconv.Atoi(responseStatusStatusThresholdStr)
	if err != nil {
		return nil, errors.New("FRONTMAN_DETECT_RESPONSE_STATUS_STATUS_THRESHOLD must be an integer")
	}

	sqlInjectionAlertThresholdStr, found := os.LookupEnv("FRONTMAN_DETECT_SQLINJECTION_ALERT_THRESHOLD")
	if !found {
		sqlInjectionAlertThresholdStr = "7" // require a score of 7 to alert by default
	}
	sqlInjectionAlertThreshold, err := strconv.Atoi(sqlInjectionAlertThresholdStr)
	if err != nil {
		return nil, errors.New("FRONTMAN_DETECT_SQLINJECTION_ALERT_THRESHOLD must be an integer")
	}

	return &ProxyConfig{
		DetectBruteForcePath:                bruteForcePath,
		DetectBruteForceUsernameField:       bruteForceUsernameField,
		DetectBruteForcePasswordField:       bruteForcePasswordField,
		DetectBruteForceAlarmThreshold:      bruteForceAlarmThreshold,
		DetectBruteForceExpireSeconds:       bruteForceExpireSeconds,
		DetectBruteForceSalt:                bruteForceSalt,
		DetectResponseStatusExpireSeconds:   responseStatusExpireSeconds,
		DetectResponseStatusIPThreshold:     responseStatusIPThreshold,
		DetectResponseStatusStatusThreshold: responseStatusStatusThreshold,
		DetectSQLInjectionAlertThreshold:    sqlInjectionAlertThreshold,
		ProxyAddress:                        proxyAddress,
		ProxyPort:                           proxyPort,
		PrivateURL:                          privateURL,
		PublicURL:                           publicURL,
		RedisHost:                           redisHost,
		RedisPassword:                       redisPassword,
		RedisDB:                             redisDB,
	}, nil
}
