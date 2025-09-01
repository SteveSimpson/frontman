package main

import (
	"context"
	"log"
	"net/http"

	"github.com/SteveSimpson/frontman/config"
	"github.com/SteveSimpson/frontman/proxy"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Ping the Redis server to check the connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Redis connection successful: ", pong)

	p, err := proxy.NewProxy(cfg, client)
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
