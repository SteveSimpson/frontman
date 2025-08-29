package config

import (
	"encoding/json"
	"os"
)

type ProxyConfig struct {
	ProxyAddress string `json:"proxy_address"`
	ProxyPort    int    `json:"proxy_port"`
	PrivateURL   string `json:"private_url"`
	PublicURL    string `json:"public_url"`
}

func LoadConfig(filePath string) (*ProxyConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &ProxyConfig{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
