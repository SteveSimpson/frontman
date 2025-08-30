package config

import (
	"errors"
	"os"
	"strconv"
)

type ProxyConfig struct {
	DetectBruteForcePath string `json:"detect_brute_force_path"`
	ProxyAddress         string `json:"proxy_address"`
	ProxyPort            string `json:"proxy_port"`
	PrivateURL           string `json:"private_url"`
	PublicURL            string `json:"public_url"`
	RedisHost            string `json:"redis_host"`
	RedisPassword        string `json:"redis_password"`
	RedisDB              int    `json:"redis_db"`
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

	return &ProxyConfig{
		DetectBruteForcePath: bruteForcePath,
		ProxyAddress:         proxyAddress,
		ProxyPort:            proxyPort,
		PrivateURL:           privateURL,
		PublicURL:            publicURL,
		RedisHost:            redisHost,
		RedisPassword:        redisPassword,
		RedisDB:              redisDB,
	}, nil
}
