package routes

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type ReverseProxy struct {
	proxy *httputil.ReverseProxy
}

func NewReverseProxy(targetHost, fromPrefix, toPrefix string) http.Handler {
	target, err := url.Parse(targetHost)
	if err != nil {
		log.Fatalf("Invalid target host: %v", err)
	}

	proxy := &httputil.ReverseProxy{
		Transport: &http.Transport{
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Director: func(req *http.Request) {
			// Store original for logging
			originalURL := *req.URL

			// Path rewriting
			if strings.HasPrefix(req.URL.Path, fromPrefix) {
				req.URL.Path = toPrefix + strings.TrimPrefix(req.URL.Path, fromPrefix)
			}

			// Update request properties
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host // Important for virtual hosting

			// Remove headers that might conflict
			req.Header.Del("X-Forwarded-Host")
			req.Header.Set("X-Forwarded-For", req.RemoteAddr)

			// Debug logging
			if originalURL.Path != req.URL.Path {
				log.Printf("Rewrote %s%s → %s%s",
					originalURL.Host, originalURL.Path,
					req.URL.Host, req.URL.Path)
			}
		},
	}

	// Response modifier to clean duplicate headers
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Remove conflicting CORS headers from downstream
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Credentials")
		resp.Header.Del("Access-Control-Expose-Headers")

		// Ensure no cache for API responses
		if strings.HasPrefix(resp.Request.URL.Path, "/api/") {
			resp.Header.Set("Cache-Control", "no-store, max-age=0")
		}

		return nil
	}

	// Enhanced error handling
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf(" Proxy error for %s %s: %v", r.Method, r.URL.Path, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{
			"error": "service_unavailable",
			"message": "Backend service not responding"
		}`))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("⇨ %s %s (%v)",
				r.Method, r.URL.Path, time.Since(start))
		}()

		// Preflight request short-circuit
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		proxy.ServeHTTP(w, r)
	})
}
