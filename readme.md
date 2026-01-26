
# rHttpbin

`rHttpbin` is a small, focused **httpbin-style** service built on top of **`rhttp`**, a custom HTTP server implemented directly over TCP in Go.

It exists for the same reason the best infrastructure tools exist: to make the invisible visible. This project is a practical, readable playground for understanding how HTTP requests are parsed, routed, inspected, and turned into responses — without hiding behind `net/http`.

The result is a fast, dependency-light demo server that ships a Swagger UI homepage and a handful of endpoints that help you test clients, debug integrations, and explore HTTP behavior.

## Highlights

- Custom HTTP server over `net.Listener` + goroutines (no `net/http`).
- Trie-based router with **static**, **parameter** routes (`:name`), and **wildcard** routes (`*`).
- Request inspection helpers (headers, IP, user-agent, URL reconstruction).
- httpbin-like “method echo” endpoints that reflect request details.
- Dynamic data endpoints (uuid, random bytes, delay, base64 decode).
- Image endpoints that serve real binary assets based on `Accept`.
- File responses (buffered) and chunked responses (supported by the framework).

## Quick start

### Prerequisites

- Go `1.25.3` (as declared in `go.mod`).

### Run

```bash
go run ./cmd/rhttpbin
```

Then open:

- Swagger UI: `http://localhost:8000/`
- OpenAPI spec: `http://localhost:8000/static/openapi.yaml`
- Health check: `http://localhost:8000/health`

The port is currently defined in [cmd/rhttpbin/main.go](cmd/rhttpbin/main.go) as `PORT`.

## API surface (httpbin-style)

This server exposes a compact set of endpoints. The most “discoverable” way to browse them is the Swagger UI at `/`, backed by the OpenAPI document served from `/static/openapi.yaml`.

### HTTP methods

- `GET  /get`
- `POST /post`
- `PUT  /put`
- `PATCH /patch`
- `DELETE /delete`

These return JSON describing the request (query args, headers, origin IP, full URL; plus body/form/files for write methods).

### Request inspection

- `GET /headers` – echoes request headers as JSON
- `GET /ip` – returns `{ "origin": "..." }`
- `GET /user-agent` – returns `{ "user-agent": "..." }`

### Response inspection

- `GET /cache` – returns `304` if `If-Modified-Since` or `If-None-Match` is present; otherwise a normal JSON payload
- `GET /cache/:value` – sets `Cache-Control: public, max-age=:value`
- `GET /etag/:etag` – supports `If-None-Match` (304) and `If-Match` (412)
- `GET|POST /response-headers` – sets response headers based on query params

### Dynamic data

- `GET /base64/:value` – decodes base64 to plain text
- `GET /uuid` – returns a UUID v4
- `GET /bytes/:n` – returns up to 102400 random bytes
- `GET|POST|PUT|PATCH|DELETE /delay/:delay` – sleeps up to 10 seconds then responds

### Images

- `GET /image` – chooses an image based on the request `Accept` header
- `GET /image/png`
- `GET /image/jpeg`
- `GET /image/webp`
- `GET /image/svg`
- `GET /image/gif`

## Curl examples

### Basic method echo

```bash
curl -s 'http://localhost:8000/get?hello=world' | jq
```

### POST JSON

```bash
curl -s \
	-H 'Content-Type: application/json' \
	-d '{"name":"rhttp"}' \
	http://localhost:8000/post | jq
```

### POST form fields

```bash
curl -s \
	-H 'Content-Type: application/x-www-form-urlencoded' \
	-d 'a=1&b=two' \
	http://localhost:8000/post | jq
```

### Multipart upload

```bash
curl -s \
	-F 'note=hello' \
	-F 'file=@./web/static/openapi.yaml' \
	http://localhost:8000/post | jq
```

### Cache behaviors

```bash
curl -i http://localhost:8000/cache
curl -i -H 'If-None-Match: anything' http://localhost:8000/cache
curl -i http://localhost:8000/cache/60
```

### ETag behaviors

```bash
curl -i http://localhost:8000/etag/mytag
curl -i -H 'If-None-Match: mytag' http://localhost:8000/etag/mytag
curl -i -H 'If-Match: different' http://localhost:8000/etag/mytag
```

### Dynamic bytes and delay

```bash
curl -i http://localhost:8000/bytes/16
curl -i http://localhost:8000/delay/2
```

### Images

```bash
curl -i -H 'Accept: image/jpeg' http://localhost:8000/image --output out.jpg
curl -i http://localhost:8000/image/png --output out.png
```

## How it works (architecture)

At a high level:

1. [cmd/rhttpbin/main.go](cmd/rhttpbin/main.go) starts the application on port `8000`.
2. [internal/application/app.go](internal/application/app.go) creates a router and starts the TCP server.
3. [internal/application/routes.go](internal/application/routes.go) wires **pages** and **API** routes.
4. `pkg/rhttp` parses the request, finds a handler in the router trie, runs it safely (panic → 500), and serializes the response.

### Router

The router in `pkg/rhttp` is a trie of URL path segments:

- Static segment: `/image/png`
- Param segment: `/cache/:value` (captured as `value`)
- Wildcard segment: `/static/*` (captured as `*` for the remaining path)

Query parameters are parsed from the request target and exposed via `QParam()` / `QParamForEach()`.

### Request parsing

`pkg/rhttp` implements a small state machine:

- Request line (`METHOD target HTTP/1.1`)
- Headers (`Name: Value`)
- Optional body based on `Content-Length`

It currently expects classic HTTP/1.1 framing using `Content-Length` (chunked *request* bodies are not implemented).

### Responses

`pkg/rhttp` supports:

- Regular responses with a fixed `Content-Length`
- Chunked responses via `Transfer-Encoding: chunked` with trailers (`x-content-sha256`, `x-content-length`)
- Convenience helpers like `ResponseJSON()` and `ErrorResponseJSON()`

## Repo layout

- [cmd/rhttpbin/main.go](cmd/rhttpbin/main.go): binary entrypoint
- [internal/application](internal/application): route wiring and app start
- [internal/application/apis](internal/application/apis): httpbin-style API handlers
- [internal/application/pages](internal/application/pages): homepage (`/`), health, and static file routing
- [internal/build](internal/build): builds the “echo” payloads (read/write)
- [internal/inspect](internal/inspect): client IP + full URL reconstruction
- [internal/respond](internal/respond): file responses for static assets
- [pkg/rhttp](pkg/rhttp): the core server, router, request/response types
- [web/static](web/static): OpenAPI spec + static assets (images, etc.)
- [web/templates](web/templates): Swagger UI HTML

## Development

Run tests:

```bash
go test ./...
```

Format code:

```bash
gofmt -w .
```

## Notes and current limitations

This is intentionally a learning-focused implementation, not a drop-in replacement for production HTTP servers:

- One request is handled per accepted TCP connection (even though the default response header advertises `Connection: keep-alive`).
- Request bodies are read using `Content-Length` only (no chunked request decoding).
- No HTTP/2, no TLS termination, no compression.
- Target validation is intentionally strict (only a limited character set is accepted for registered routes).

## Contact

- Author: MSI Rony
- Email: `abrony@gmail.com`
- Website: <https://sirony.site>

