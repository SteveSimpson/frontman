package main

import (
	"log"
	"net/http"

	"github.com/SteveSimpson/frontman/config"
	"github.com/SteveSimpson/frontman/proxy"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	p, err := proxy.NewProxy(cfg)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	http.HandleFunc("/", p.HandleRequest)

	log.Printf("Starting web proxy server on %s:%s\n", cfg.ProxyAddress, cfg.ProxyPort)
	err = http.ListenAndServe(cfg.ProxyAddress+":"+cfg.ProxyPort, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
