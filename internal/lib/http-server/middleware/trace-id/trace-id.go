package trace_id

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const (
	CtxTraceId    string = "traceId"
	headerTraceId string = "trace-id"
)

func GetTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(CtxTraceId).(string); ok {
		return v
	}

	return ""
}

func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceId := r.Header.Get(headerTraceId)

			if traceId == "" {
				traceId = randomHex(16)

				w.Header().Set(headerTraceId, traceId)
			}

			ctx := context.WithValue(r.Context(), CtxTraceId, traceId)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func randomHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}
