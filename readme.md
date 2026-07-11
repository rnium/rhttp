# rhttp

`rhttp` is a minimal HTTP/1.1 server framework built directly over TCP in Go, with no dependency on `net/http`. It implements a focused subset of [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110) (HTTP Semantics) and [RFC 9112](https://www.rfc-editor.org/rfc/rfc9112) (HTTP/1.1 Message Syntax).

The goal is clarity over completeness — every layer from TCP connection acceptance to response serialization is explicit and readable.

## Features

- Raw TCP server via `net.Listener` with one goroutine per connection
- Keep-alive support with a 60-second idle timeout
- Graceful shutdown with `sync.WaitGroup` and a quit signal channel
- Trie-based router with **static**, **parameter** (`:name`), and **wildcard** (`*`) segments
- Hand-rolled request parser — state machine over a pooled read buffer (request line → headers → body)
- Request body reading via `Content-Length` (chunked request bodies not supported)
- Chunked responses via `Transfer-Encoding: chunked` with `x-content-sha256` and `x-content-length` trailers
- RFC 9110 §5.6.2 compliant header token validation; duplicate headers are combined per §5.3
- Cookie serialization (`Set-Cookie`) and cookie parsing from the `Cookie` request header
- URL query parameter parsing with percent-decoding
- Multipart form, URL-encoded form, and raw body parsing
- Panic recovery in handlers — panics are caught and turned into `500` responses
- Response helpers: `ResponseJSON`, `ErrorResponseJSON`, `Redirect`, `NewChunkedResponse`

## Quick start

### Prerequisites

Go `1.25.3` or later.

### Installation

```bash
go get github.com/rnium/rhttp
```

### Minimal example

```go
package main

import (
    "github.com/rnium/rhttp/pkg/rhttp"
)

func main() {
    router := rhttp.NewRouter()

    router.Get("/hello/:name", func(req *rhttp.Request) *rhttp.Response {
        name, _ := req.Param("name")
        return rhttp.ResponseJSON(200, map[string]string{"hello": name})
    })

    server := rhttp.NewServer(router)
    server.Start(8080)

    select {} // block forever; call server.Close() to shut down
}
```

## Core concepts

### Server

`NewServer(router)` creates a server. `Start(port)` begins accepting TCP connections; each connection is handled in its own goroutine. Call `Close()` to stop the listener and wait for in-flight connections to finish.

### Router

The router is a trie of URL path segments. Register handlers with the HTTP-method helpers:

```go
router.Get("/static/path", handler)
router.Post("/users/:id", handler)
router.Delete("/files/*", handler)
```

Inside a handler, read captured segments with `req.Param("id")` and query parameters with `req.QParam("key")` or iterate them with `req.QParamForEach(fn)`.

### Request

The `Request` type exposes:

| Field / Method | Description |
|---|---|
| `RequestLine` | `Method`, `Target`, `Version` |
| `Headers` | Case-insensitive header store |
| `Body` | Raw body bytes (populated when `Content-Length > 0`) |
| `Param(name)` | URL path parameter |
| `QParam(name)` | Query parameter |
| `FormData()` | Parsed form — multipart, URL-encoded, or raw |
| `Cookies()` | Map of request cookies |
| `RemoteAddr()` | Client's TCP address |

### Response

Create a fixed-length response or a streaming chunked response:

```go
// Fixed-length
res := rhttp.NewResponse(200, []byte("hello"))
res.SetHeader("Content-Type", "text/plain")

// Chunked (streams from any io.Reader)
res := rhttp.NewChunkedResponse(200, file)

// Convenience helpers
res := rhttp.ResponseJSON(200, payload)
res := rhttp.ErrorResponseJSON(400, "bad input")
res := rhttp.Redirect("/new-path", false)
```

Chunked responses automatically compute and send `x-content-sha256` and `x-content-length` trailers. The `io.Reader` is closed after the response is written if it implements `io.Closer`.

### Headers

Response headers are set with `res.SetHeader(name, value)`. A set of headers that control the HTTP framing (`Content-Length`, `Transfer-Encoding`, `Connection`, etc.) are protected and cannot be overwritten by handlers.

### Cookies

```go
// Set a cookie on the response
res.SetCookie(&rhttp.Cookie{
    Name:     "session",
    Value:    "abc123",
    HttpOnly: true,
    Secure:   true,
    MaxAge:   "3600",
})

// Read cookies from the request
cookies := req.Cookies() // map[string]string
```

## Repo layout

```
cmd/rhttpbin/          – rHttpbin binary entrypoint
internal/              – application-level code for rHttpbin
pkg/rhttp/             – the framework (server, router, request, response, headers, cookies, forms)
web/                   – static assets and templates for rHttpbin
```

## Running tests

```bash
go test ./...
```

## rHttpbin

[`rHttpbin`](https://rhttpbin.sirony.site/) is a live httpbin-style HTTP inspection service built on top of `rhttp`. It demonstrates the framework in a real workload and serves as an integration test for its features.

Browse the full API surface through the **Swagger UI** at [rhttpbin.sirony.site](https://rhttpbin.sirony.site/) or inspect the raw OpenAPI spec at `/static/openapi.yaml`.

## Known limitations

- One HTTP/1.1 request–response cycle per TCP connection per iteration (keep-alive is negotiated but pipelining is not supported).
- Request bodies are read by `Content-Length` only — chunked request transfer encoding is not decoded.
- No HTTP/2, no TLS, no compression.
- Route target characters are restricted to letters, digits, and `*./-_:`.
