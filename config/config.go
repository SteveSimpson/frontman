package config

import (
	"errors"
	"os"
)

type ProxyConfig struct {
	ProxyAddress string `json:"proxy_address"`
	ProxyPort    string `json:"proxy_port"`
	PrivateURL   string `json:"private_url"`
	PublicURL    string `json:"public_url"`
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

	return &ProxyConfig{
		ProxyAddress: proxyAddress,
		ProxyPort:    proxyPort,
		PrivateURL:   privateURL,
		PublicURL:    publicURL,
	}, nil
}
