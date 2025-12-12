# HTTP Tracing Middleware

OpenTelemetry HTTP middleware for Go with support for multiple frameworks.

## Features

- ✅ **OpenTelemetry Distributed Tracing** - Automatic span creation for HTTP requests
- ✅ **Context Propagation** - Extract and inject trace context from/to headers
- ✅ **Request ID Generation** - Automatic request ID with `X-Request-ID` header
- ✅ **Status Code Tracking** - Capture response status and mark errors
- ✅ **Error Recording** - Record panics and errors in spans
- ✅ **Structured Logging** - Context-aware logs with trace IDs
- ✅ **Multiple Frameworks** - stdlib, Fiber (Gin coming soon)
- ✅ **Configurable** - Skip paths, disable logging, customize behavior

## Installation

```bash
go get github.com/ubin/go-telemetry
```

## Usage

### Standard net/http

```go
import (
    "net/http"
    httpMiddleware "github.com/ubin/go-telemetry/middleware/http"
    "github.com/ubin/go-telemetry/telemetry"
)

// Initialize telemetry
tp, err := telemetry.InitTracer(cfg)
if err != nil {
    log.Fatal(err)
}

// Configure middleware
middlewareConfig := &httpMiddleware.Config{
    TracerProvider:    tp,
    Logger:            logger,
    ServiceName:       "my-service",
    SkipPaths:         []string{"/health", "/metrics"},
    GenerateRequestID: true,
}

// Wrap your handler
mux := http.NewServeMux()
mux.HandleFunc("/", homeHandler)

handler := httpMiddleware.Middleware(middlewareConfig)(mux)

http.ListenAndServe(":8080", handler)
```

### Fiber

```go
import (
    "github.com/gofiber/fiber/v2"
    httpMiddleware "github.com/ubin/go-telemetry/middleware/http"
)

app := fiber.New()

// Add tracing middleware
app.Use(httpMiddleware.FiberMiddleware(&httpMiddleware.Config{
    TracerProvider:    tp,
    Logger:            logger,
    ServiceName:       "my-service",
    SkipPaths:         []string{"/health"},
    GenerateRequestID: true,
}))

app.Get("/", func(c *fiber.Ctx) error {
    // Access trace info from context
    traceID := c.Locals("trace_id").(string)
    requestID := c.Locals("request_id").(string)

    return c.JSON(fiber.Map{
        "trace_id": traceID,
        "request_id": requestID,
    })
})

app.Listen(":8080")
```

## Configuration

```go
type Config struct {
    // TracerProvider is the OpenTelemetry tracer provider
    // If nil, tracing will be skipped (but logging still works)
    TracerProvider *trace.TracerProvider

    // Logger for request logging (optional)
    // Accepts any logger that implements logger.ContextLogger interface
    // (InfoContext and ErrorContext methods)
    Logger logger.ContextLogger

    // ServiceName for tracing spans (default: "http-server")
    ServiceName string

    // SkipPaths to exclude from tracing (e.g., /health, /metrics)
    SkipPaths []string

    // SkipLogging disables request logging
    SkipLogging bool

    // GenerateRequestID enables X-Request-ID header (default: true)
    GenerateRequestID bool
}
```

## Response Headers

The middleware automatically adds trace context to response headers:

- `X-Request-ID` - Unique request identifier
- `X-Trace-ID` - OpenTelemetry trace ID
- `X-Span-ID` - OpenTelemetry span ID

## Span Attributes

Each HTTP request span includes:

- `http.method` - HTTP method (GET, POST, etc.)
- `http.path` - Request path
- `http.route` - Route pattern (e.g., `/users/:id`)
- `http.scheme` - Protocol (http/https)
- `http.target` - Full request URI
- `http.host` - Host header
- `http.user_agent` - User agent string
- `http.remote_addr` - Client IP address
- `http.status_code` - Response status code
- `http.response_size` - Response body size in bytes
- `http.request_id` - Request ID (if enabled)

## Error Handling

- **Status >= 500**: Span marked as error with `codes.Error`
- **Panics**: Automatically recovered, recorded in span, and re-panicked
- **Fiber errors**: Recorded using `span.RecordError()`

## Logging

Logs include trace context automatically when using `*Context` methods:

```go
logger.InfoContext(ctx, "Processing request", "user_id", 123)
// Output: {"msg":"Processing request","trace_id":"abc123","span_id":"def456","user_id":123}
```

## Creating Child Spans

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Create child span for database operation
    tracer := otel.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "database-query")
    defer span.End()

    // Logs will include parent trace ID
    logger.InfoContext(ctx, "Querying database")

    // Do database work...
}
```

## Example: Complete Setup

See the [example directory](../../example) for a complete working example.

## License

MIT
