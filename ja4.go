package ja4

import (
	"net/http"

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
type JA4Placeholder struct{}

// CaddyModule returns the Caddy module information.
func (JA4Placeholder) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.ja4_placeholder",
		New: func() caddy.Module { return new(JA4Placeholder) },
	}
}

// Provision sets up the module.
func (p *JA4Placeholder) Provision(ctx caddy.Context) error {
	return nil
}

// Validate ensures the module is properly configured. (No-op)
func (p *JA4Placeholder) Validate() error {
	return nil
}

// ServeHTTP calculates the JA4 hash and sets it as a placeholder, then passes through to the next handler.
func (p JA4Placeholder) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Calculate JA4 hash and set it as a placeholder
	ja4Hash := calculateJA4(r)

	// Get the replacer from the request context and set the placeholder
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	repl.Set("ja4h", ja4Hash)

	// Pass through to the next handler
	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (p *JA4Placeholder) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	// This directive takes no arguments
	d.Next() // consume directive name
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new JA4Placeholder.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var p JA4Placeholder
	err := p.UnmarshalCaddyfile(h.Dispenser)
	return p, err
}

// calculateJA4 computes a basic JA4 hash from request properties.
func calculateJA4(r *http.Request) string {
	// Use go-ja4h to calculate the JA4 hash from the request
	hash := ja4h.JA4H(r)
	return hash
}

// Interface guards
var (
	_ caddy.Module                = (*JA4Placeholder)(nil)
	_ caddyhttp.MiddlewareHandler = (*JA4Placeholder)(nil)
	_ caddy.Provisioner           = (*JA4Placeholder)(nil)
	_ caddyfile.Unmarshaler       = (*JA4Placeholder)(nil)
)
