package logger

import (
	"log/slog"
	"net/http"
	"plata_currency_quotation/internal/api"
	"plata_currency_quotation/internal/lib/http-server/middleware/trace-id"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, api.SwaggerEndpoint) {
				next.ServeHTTP(w, r)

				return
			}

			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("trace_id", trace_id.GetTraceID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			next.ServeHTTP(ww, r)

			entry.Info("request completed",
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.String("duration", time.Since(t1).String()),
			)
		}

		return http.HandlerFunc(fn)
	}
}
