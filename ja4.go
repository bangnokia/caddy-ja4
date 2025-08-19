package ja4

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/lum8rjack/go-ja4h/ja4h"
)

func init() {
	caddy.RegisterModule(JA4Placeholder{})
}

// JA4Placeholder implements a Caddy module that provides JA4 hash as a placeholder.
type JA4Placeholder struct{}

// CaddyModule returns the Caddy module information.
func (JA4Placeholder) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.ja4_placeholder",
		New: func() caddy.Module { return new(JA4Placeholder) },
	}
}

// Provision sets up the module and registers the placeholder.
func (p *JA4Placeholder) Provision(ctx caddy.Context) error {
	// Register the {ja4h} placeholder
	caddyhttp.RegisterPlaceholder("ja4h", func(r *http.Request, repl *caddy.Replacer) (interface{}, bool) {
		ja4Hash := calculateJA4(r)
		return ja4Hash, true
	})
	return nil
}

// Validate ensures the module is properly configured. (No-op)
func (p *JA4Placeholder) Validate() error {
	return nil
}

// ServeHTTP is a no-op handler since we only provide placeholders.
func (p JA4Placeholder) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Just pass through to the next handler
	return next.ServeHTTP(w, r)
}

// calculateJA4 computes a basic JA4 hash from request properties.
func calculateJA4(r *http.Request) string {
	// Use go-ja4h to calculate the JA4 hash from the request
	hash := ja4h.Calculate(r)
	return hash
}

// Interface guards
var (
	_ caddy.Module                = (*JA4Placeholder)(nil)
	_ caddyhttp.MiddlewareHandler = (*JA4Placeholder)(nil)
	_ caddy.Provisioner           = (*JA4Placeholder)(nil)
)
