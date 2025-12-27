// Package gateway provides the HTTP API gateway server
package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Shared HTTP client for backend requests
// http.Client is safe for concurrent use by multiple goroutines
var defaultHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

// Server implements the API gateway HTTP server
type Server struct {
	addr              string
	httpServer        *http.Server
	router            *Router
	rateLimiters      *MultiCircuitBreaker
	userRateLimiters  *UserRateLimiter
	logger            *log.Logger
	middleware        []Middleware
	shutdownOnce      sync.Once
	shutdownTimeout   time.Duration
}

// ServerConfig holds configuration for the gateway
type ServerConfig struct {
	Address              string
	Port                 int
	ShutdownTimeout      time.Duration
	RateLimitConfig      CircuitBreakerConfig
	DefaultRateLimitRPS  int
	MaxConcurrentPerUser int
}

// Middleware is a function that wraps an HTTP handler
type Middleware func(http.Handler) http.Handler

// NewServer creates a new gateway server
func NewServer(config ServerConfig) *Server {
	if config.ShutdownTimeout == 0 {
		config.ShutdownTimeout = 30 * time.Second
	}

	addr := fmt.Sprintf("%s:%d", config.Address, config.Port)

	s := &Server{
		addr:             addr,
		router:           NewRouter(),
		logger:           log.New(io.Discard, "[Gateway] ", log.LstdFlags),
		middleware:       make([]Middleware, 0),
		shutdownTimeout:  config.ShutdownTimeout,
		rateLimiters:     NewMultiCircuitBreaker(config.RateLimitConfig),
		userRateLimiters: NewUserRateLimiter(config.DefaultRateLimitRPS, 1.0, 5*time.Minute),
	}

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         s.addr,
		Handler:      s.handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

// RegisterRoute registers a route with a backend service
func (s *Server) RegisterRoute(pattern string, backendURL string) error {
	backendParsed, err := url.Parse(backendURL)
	if err != nil {
		return fmt.Errorf("invalid backend URL: %w", err)
	}

	s.router.RegisterRoute(pattern, backendURL, NewBackendHandler(backendParsed))
	s.logger.Printf("Registered route: %s -> %s", pattern, backendURL)
	return nil
}

// AddMiddleware adds middleware to the chain
func (s *Server) AddMiddleware(m Middleware) {
	s.middleware = append(s.middleware, m)
}

// handler creates the HTTP handler with middleware chain
func (s *Server) handler() http.Handler {
	baseHandler := http.Handler(http.HandlerFunc(s.handleRequest))

	// Apply middleware in reverse order (last added executes first)
	handler := baseHandler
	for i := len(s.middleware) - 1; i >= 0; i-- {
		handler = s.middleware[i](handler)
	}

	return handler
}

// handleRequest handles incoming HTTP requests
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Extract user ID from header (if present)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	// Check rate limiter for user
	if !s.userRateLimiters.Allow(userID, 1) {
		s.respondError(w, http.StatusTooManyRequests, "Rate limit exceeded")
		s.logger.Printf("Rate limit exceeded for user: %s", userID)
		return
	}

	// Check circuit breaker for endpoint
	route := s.router.FindRoute(r.Method, r.URL.Path)
	if route == nil {
		s.respondError(w, http.StatusNotFound, "Route not found")
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), route.Timeout)
	defer cancel()
	r = r.WithContext(ctx)

	// Call through circuit breaker
	var backendResp *http.Response
	var backendErr error

	err := s.rateLimiters.Call(route.Pattern, func() error {
		backendResp, backendErr = route.Handler(r)
		return backendErr
	})

	// Handle circuit breaker errors
	if err == ErrCircuitOpen {
		s.respondError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
		s.logger.Printf("Circuit breaker OPEN for route: %s", route.Pattern)
		return
	}

	// Handle backend errors
	if backendErr != nil {
		s.respondError(w, http.StatusBadGateway, "Backend error")
		s.logger.Printf("Backend error: %v", backendErr)
		return
	}

	if backendResp == nil {
		s.respondError(w, http.StatusInternalServerError, "No response from backend")
		return
	}

	// Copy response headers
	for k, v := range backendResp.Header {
		w.Header()[k] = v
	}

	// Copy status code
	w.WriteHeader(backendResp.StatusCode)

	// Copy body
	if _, err := io.Copy(w, backendResp.Body); err != nil {
		s.logger.Printf("Error copying response body: %v", err)
	}
	backendResp.Body.Close()

	// Log request
	duration := time.Since(start)
	s.logger.Printf(
		"%s %s - User: %s - Status: %d - Duration: %v",
		r.Method,
		r.URL.Path,
		userID,
		backendResp.StatusCode,
		duration,
	)
}

// respondError sends an error response
func (s *Server) respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"error":   message,
		"status":  statusCode,
		"time":    time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// Start starts the gateway server
func (s *Server) Start() error {
	s.logger.Printf("Starting gateway on %s", s.addr)
	return s.httpServer.ListenAndServe()
}

// StartAsync starts the server in a goroutine
func (s *Server) StartAsync() <-chan error {
	errChan := make(chan error, 1)
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()
	return errChan
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	var shutdownErr error
	s.shutdownOnce.Do(func() {
		s.logger.Println("Shutting down gateway")

		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("shutdown error: %w", err)
		}

		s.logger.Println("Gateway shutdown complete")
	})
	return shutdownErr
}

// BackendHandler handles proxying to a backend service
type BackendHandler struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

// NewBackendHandler creates a backend handler
func NewBackendHandler(target *url.URL) func(*http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		// Clone the request to avoid modifying the original
		backendReq := req.Clone(req.Context())
		
		// Update URL for backend
		backendReq.URL.Scheme = target.Scheme
		backendReq.URL.Host = target.Host
		backendReq.Host = target.Host
		backendReq.RequestURI = "" // Must be empty for client requests

		// Add forwarding headers
		backendReq.Header.Add("X-Forwarded-For", getClientIP(req))
		backendReq.Header.Add("X-Forwarded-Proto", req.Proto)
		backendReq.Header.Add("X-Forwarded-Host", req.Host)

		// Use shared HTTP client (safe for concurrent use)
		return defaultHTTPClient.Do(backendReq)
	}
}

// Route represents a registered route
type Route struct {
	Pattern string
	Backend string
	Timeout time.Duration
	Handler func(*http.Request) (*http.Response, error)
}

// Router manages routes
type Router struct {
	routes map[string]*Route
	mu     sync.RWMutex
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]*Route),
	}
}

// RegisterRoute registers a route
func (r *Router) RegisterRoute(pattern string, backend string, handler func(*http.Request) (*http.Response, error)) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routes[pattern] = &Route{
		Pattern: pattern,
		Backend: backend,
		Timeout: 10 * time.Second,
		Handler: handler,
	}
}

// FindRoute finds a route matching the request
func (r *Router) FindRoute(method, path string) *Route {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Simple prefix matching (more sophisticated routing could be implemented)
	for pattern, route := range r.routes {
		if strings.HasPrefix(path, pattern) {
			return route
		}
	}

	return nil
}

// Helper functions

func copyHeaders(h http.Header) http.Header {
	header := make(http.Header)
	for k, v := range h {
		header[k] = append([]string(nil), v...)
	}
	return header
}

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	if forwarded := r.Header.Get("X-Real-IP"); forwarded != "" {
		return forwarded
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

// LoggingMiddleware adds logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// AuthMiddleware adds authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for auth header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Simple token validation (in production, use real validation)
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
