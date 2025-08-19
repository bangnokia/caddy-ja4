package ja4

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/lum8rjack/go-ja4h"
)

func init() {
	caddy.RegisterModule(JA4Placeholder{})
	httpcaddyfile.RegisterHandlerDirective("ja4_placeholder", parseCaddyfile)
}

// JA4Placeholder implements a Caddy module that provides JA4 hash as a placeholder.
type JA4Placeholder struct {
	// Cache configuration
	CacheDuration string `json:"cache_duration,omitempty"`

	// Internal cache
	cache     map[string]cacheEntry
	cacheMu   sync.RWMutex
	cacheTTL  time.Duration
}

type cacheEntry struct {
	hash      string
	timestamp time.Time
}

// CaddyModule returns the Caddy module information.
func (JA4Placeholder) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.ja4_placeholder",
		New: func() caddy.Module { return new(JA4Placeholder) },
	}
}

// Provision sets up the module.
func (p *JA4Placeholder) Provision(ctx caddy.Context) error {
	// Set default cache duration
	if p.CacheDuration == "" {
		p.CacheDuration = "30s"
	}

	// Parse cache duration
	var err error
	p.cacheTTL, err = time.ParseDuration(p.CacheDuration)
	if err != nil {
		p.cacheTTL = 30 * time.Second
	}

	// Initialize cache
	p.cache = make(map[string]cacheEntry)

	return nil
}

// Validate ensures the module is properly configured. (No-op)
func (p *JA4Placeholder) Validate() error {
	return nil
}

// ServeHTTP calculates the JA4 hash and sets it as a placeholder, then passes through to the next handler.
func (p *JA4Placeholder) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Get client identifier (IP + User-Agent for better caching)
	clientKey := r.RemoteAddr + "|" + r.UserAgent()

	// Check cache first
	ja4Hash := p.getCachedHash(clientKey)
	if ja4Hash == "" {
		// Calculate and cache JA4 hash
		ja4Hash = calculateJA4(r)
		p.setCachedHash(clientKey, ja4Hash)
	}

	// Get or create the replacer from the request context
	repl, ok := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	if !ok || repl == nil {
		repl = caddy.NewReplacer()
		ctx := context.WithValue(r.Context(), caddy.ReplacerCtxKey, repl)
		r = r.WithContext(ctx)
	}
	repl.Set("ja4h", ja4Hash)

	// Pass through to the next handler
	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (p *JA4Placeholder) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume directive name

	// Optional: parse cache duration argument
	if d.NextArg() {
		p.CacheDuration = d.Val()
	}

	return nil
}

// parseCaddyfile unmarshals tokens from h into a new JA4Placeholder.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var p JA4Placeholder
	err := p.UnmarshalCaddyfile(h.Dispenser)
	return &p, err
}

// getCachedHash retrieves a cached JA4 hash if it exists and is not expired
func (p *JA4Placeholder) getCachedHash(clientKey string) string {
	p.cacheMu.RLock()
	defer p.cacheMu.RUnlock()

	entry, exists := p.cache[clientKey]
	if !exists {
		return ""
	}

	// Check if cache entry is expired
	if time.Since(entry.timestamp) > p.cacheTTL {
		return ""
	}

	return entry.hash
}

// setCachedHash stores a JA4 hash in the cache
func (p *JA4Placeholder) setCachedHash(clientKey, hash string) {
	p.cacheMu.Lock()
	defer p.cacheMu.Unlock()

	// Clean up expired entries periodically
	now := time.Now()
	for key, entry := range p.cache {
		if now.Sub(entry.timestamp) > p.cacheTTL {
			delete(p.cache, key)
		}
	}

	// Store new entry
	p.cache[clientKey] = cacheEntry{
		hash:      hash,
		timestamp: now,
	}
}

// calculateJA4 computes a basic JA4 hash from request properties.
func calculateJA4(r *http.Request) string {
	// Use go-ja4h to calculate the JA4 hash from the request
	hash := ja4h.JA4H(r)

	// Provide fallback for empty/invalid hashes
	if hash == "" {
		return "unknown"
	}

	return hash
}

// Interface guards
var (
	_ caddy.Module                = (*JA4Placeholder)(nil)
	_ caddyhttp.MiddlewareHandler = (*JA4Placeholder)(nil)
	_ caddy.Provisioner           = (*JA4Placeholder)(nil)
	_ caddyfile.Unmarshaler       = (*JA4Placeholder)(nil)
)
