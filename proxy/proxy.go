package proxy

import (
	"bytes"
	"io"
	"log"
	"net"
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

// responseRecorder wraps http.ResponseWriter to capture status and body
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (p *Proxy) HandleRequest(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = p.publicUrl.Scheme
	r.URL.Host = p.publicUrl.Host
	r.Host = p.publicUrl.Host

	requestClone := r.Clone(r.Context())

	clientIP := requestClone.RemoteAddr
	if host, _, err := net.SplitHostPort(clientIP); err == nil {
		clientIP = host
	}

	client := detect.ClientInfo{
		IP:         clientIP,
		Method:     r.Method,
		RequestURI: r.RequestURI,
		UserAgent:  r.UserAgent(),
		Referrer:   r.Referer(),
	}

	// read the body (this will consume it)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	// restore the io.ReadCloser to its original state
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	requestClone.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	log.Printf("REQUEST for %s from %s", r.RequestURI, r.RemoteAddr)

	// call a goroutine to check for attacks (asynchronously)
	go detect.RequestDetectors(requestClone)

	// could call a function to block requests here if needed

	// Wrap the ResponseWriter to capture status and body
	rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
	p.proxy.ServeHTTP(rr, r)

	// Clone the response body and status code
	responseBody := rr.body.Bytes()
	responseStatus := rr.statusCode

	go detect.ResponseDetectors(client, responseStatus, responseBody)
}
