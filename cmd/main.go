package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/SteveSimpson/frontman/internal/config"
	"github.com/SteveSimpson/frontman/internal/proxy"
)

func main() {
	cfg, err := config.LoadConfig("frontman_config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	p, err := proxy.NewProxy(cfg)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	http.HandleFunc("/", p.HandleRequest)

	log.Printf("Starting web proxy server on %s:%d\n", cfg.ProxyAddress, cfg.ProxyPort)
	err = http.ListenAndServe(cfg.ProxyAddress+":"+strconv.Itoa(cfg.ProxyPort), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
