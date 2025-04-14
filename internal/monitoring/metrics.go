package monitoring

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by method and path",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)

	cacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	cacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(cacheHits)
	prometheus.MustRegister(cacheMisses)
}

// MetricsMiddleware returns a middleware that collects HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		activeConnections.Inc()
		defer activeConnections.Dec()

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			string(rune(status)),
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}

// IncrementCacheHits increments the cache hits counter
func IncrementCacheHits() {
	cacheHits.Inc()
}

// IncrementCacheMisses increments the cache misses counter
func IncrementCacheMisses() {
	cacheMisses.Inc()
}

// RegisterMetricsEndpoint registers the /metrics endpoint for Prometheus scraping
func RegisterMetricsEndpoint(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
