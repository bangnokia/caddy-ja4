
# Caddy JA4 Hash Placeholder

This is a module for the Caddy web server that provides JA4 hash as a placeholder that can be injected into response bodies, logs, or anywhere placeholders are supported.

## Installation
You can build with xcaddy:
```bash
xcaddy build --with github.com/bangnokia/caddy-ja4
```

## Usage

The module provides a `{ja4h}` placeholder that you can use anywhere in your Caddyfile.

### Basic Response Example

```caddyfile
route {
    ja4_placeholder
    respond "Your JA4 hash is: {ja4h}"
}
```

### Template Example

```caddyfile
route {
    ja4_placeholder
    templates
    respond `
    <html>
        <body>
            <h1>JA4 Hash: {ja4h}</h1>
            <p>Your browser fingerprint is: {ja4h}</p>
        </body>
    </html>
    ` 200 {
        Content-Type "text/html"
    }
}
```

### JSON Response Example

```caddyfile
route {
    ja4_placeholder
    respond `{"ja4": "{ja4h}", "timestamp": "{time.now}"}` 200 {
        Content-Type "application/json"
    }
}
```

### Logging Example

```caddyfile
{
    log {
        format json
        output file /var/log/caddy/access.log
        include http.request.ja4h
    }
}

route {
    ja4_placeholdere
    respond "Hello World"
}
```

## How it works

The module registers a `{ja4h}` placeholder that calculates the JA4 hash from the incoming HTTP request. You need to include the `ja4_placeholder` handler in your route to enable the placeholder functionality.

