package metrics

import (
	"net/http"
	"plata_currency_quotation/internal/lib/metrics"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()

			next.ServeHTTP(ww, r)

			duration := time.Since(start).Seconds()

			metrics.HttpRequestsTotal.WithLabelValues(
				r.Method,
				r.URL.Path,
				strconv.Itoa(ww.Status()),
			).Inc()

			metrics.HttpRequestDuration.WithLabelValues(
				r.Method,
				r.URL.Path,
			).Observe(duration)
		})
	}
}
