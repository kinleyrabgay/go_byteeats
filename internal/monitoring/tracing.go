package monitoring

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig holds configuration for tracing
type TracingConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	JaegerEndpoint string
}

// InitTracer initializes the OpenTelemetry tracer
func InitTracer(config TracingConfig) (func(), error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.JaegerEndpoint)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
			attribute.String("environment", config.Environment),
		)),
	)

	otel.SetTracerProvider(tp)

	return func() {
		_ = tp.Shutdown(context.Background())
	}, nil
}

// TracingMiddleware returns a middleware that adds tracing to requests
func TracingMiddleware() gin.HandlerFunc {
	tracer := otel.Tracer("")

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		ctx, span := tracer.Start(ctx, path,
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.path", path),
			),
		)
		defer span.End()

		// Add trace context to response headers
		c.Writer.Header().Set("X-Trace-ID", span.SpanContext().TraceID().String())

		// Store span in context
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// Add response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)

		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("error", c.Errors.String()))
		}
	}
}
