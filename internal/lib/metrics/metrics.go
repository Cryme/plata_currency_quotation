package metrics

import (
	"log/slog"
	"net/http"
	"os"
	"plata_currency_quotation/internal/lib/logger/sl"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type SetupMetricsInterface interface {
	SetupMetrics(reg *prometheus.Registry)
}

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "incoming_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "incoming_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func Run(port uint16, services ...SetupMetricsInterface) {
	var addr string
	reg := prometheus.NewRegistry()

	log := sl.Log.With(
		slog.String("component", "middleware/metrics"),
	)

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		HttpRequestsTotal,
		HttpRequestDuration,
	)

	for _, service := range services {
		service.SetupMetrics(reg)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	addr = "localhost:" + strconv.Itoa(int(port))

	go func() {
		log.Info("metrics are available at endpoint " + addr + "/metrics")

		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Error("failed to start metrics server", sl.Err(err))

			os.Exit(1)
		}
	}()
}
