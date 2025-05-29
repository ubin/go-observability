# Go Telemetry

Go Telemetry is a flexible and modular library for integrating OpenTelemetry tracing and logging into your Go applications. It provides support for multiple tracing exporters, context-aware logging, and seamless integration with popular observability tools like Sentry.

## Features

- **OpenTelemetry Integration**: Easily initialize and configure OpenTelemetry tracing.
- **Multiple Exporters**: Supports GRPC, HTTP, Stdout, and Sentry exporters.
- **Context-Aware Logging**: Log messages with context for better traceability.
- **Stack Trace Capture**: Record errors with stack traces for debugging.
- **Modular Design**: Interfaces for basic logging, context-aware logging, and underlying logger access.
- **Easy Configuration**: Flexible configuration options for tracing and logging.


## Installation

```sh
$ go get github.com/ubin/go-telemetry
```

### Usage
Refer to the example applications in the repository for detailed usage:

- Tracing Example with Sentry: [example/sentry](example/sentry/main.go)
- Example with stdout: [examples/stdout](example/stdout/main.go)

These examples demonstrate how to configure and use tracing, logging, and error monitoring in your applications.

### Configuration

#### Tracing Configuration

The TracingConfig struct allows you to configure tracing options:
```
type TracingConfig struct {
    ServiceName       string
    Environment       string
    Enabled           bool
    ExporterType      ExporterType
    CollectorEndpoint string
    Insecure          bool
    DebugMode         bool
}
```

#### Logger Configuration
The LoggerConfig struct allows you to configure logging options:
```
type LoggerConfig struct {
    Code         string
    Level        string
    Formatter    string
    EnableCaller bool
    FileEnabled  bool
    Filename     string
    MaxSize      int
    MaxBackups   int
    MaxAge       int
    Compress     bool
    LocalTime    bool
}
```

#### Supported Exporters
- Stdout: Logs traces to the console.
- HTTP: Sends traces to an HTTP endpoint.
- GRPC: Sends traces to a GRPC endpoint.
- Sentry: Integrates with Sentry for error monitoring.

### License
MIT