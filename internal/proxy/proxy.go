package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/SteveSimpson/frontman/internal/config"
)

type Proxy struct {
	privateUrl *url.URL
	publicUrl  *url.URL
	proxy      *httputil.ReverseProxy
}

func NewProxy(cfg *config.ProxyConfig) (*Proxy, error) {
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

	return &Proxy{
		privateUrl: privateURL,
		publicUrl:  publicURL,
		proxy:      httputil.NewSingleHostReverseProxy(privateURL),
	}, nil
}

func (p *Proxy) HandleRequest(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = p.publicUrl.Scheme
	r.URL.Host = p.publicUrl.Host
	r.Host = p.publicUrl.Host

	log.Printf("Request for %s from %s", r.RequestURI, r.RemoteAddr)

	p.proxy.ServeHTTP(w, r)
}
