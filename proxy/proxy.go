package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/SteveSimpson/frontman/config"
	"github.com/SteveSimpson/frontman/detect"
	"github.com/redis/go-redis/v9"
)

type Proxy struct {
	privateUrl  *url.URL
	publicUrl   *url.URL
	proxy       *httputil.ReverseProxy
	redisClient *redis.Client
	config      *config.ProxyConfig
}

func NewProxy(cfg *config.ProxyConfig, redisClient *redis.Client) (*Proxy, error) {
	publicURL, err := url.Parse(cfg.PublicURL)
	if err != nil {
		log.Fatalf("Failed to parse PublicURL: %v", err)
		return nil, err
	}
	privateURL, err := url.Parse(cfg.PrivateURL)
	if err != nil {
		log.Fatalf("Failed to parse PrivateURL: %v", err)
		return nil, err
	}

	detect.InitDetectors(cfg, redisClient)

	return &Proxy{
		privateUrl:  privateURL,
		publicUrl:   publicURL,
		proxy:       httputil.NewSingleHostReverseProxy(privateURL),
		redisClient: redisClient,
		config:      cfg,
	}, nil
}

func (p *Proxy) HandleRequest(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = p.publicUrl.Scheme
	r.URL.Host = p.publicUrl.Host
	r.Host = p.publicUrl.Host

	clone := r.Clone(r.Context())

	// read the body (this will consume it)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	// restore the io.ReadCloser to its original state
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	clone.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// if we wanted to block the request, we could do it here

	log.Printf("Request for %s from %s", r.RequestURI, r.RemoteAddr)

	p.proxy.ServeHTTP(w, r)

	// call a goroutine to check for attacks (asynchronously)
	go detect.RunDetectors(clone)
}
