package traceparent

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
)

type CtxKeys string

const (
	CtxTraceID CtxKeys = "traceID"
	CtxSpanID  CtxKeys = "spanID"
)

func GetTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(CtxTraceID).(string); ok {
		return v
	}

	return ""
}

func GetSpanID(ctx context.Context) string {
	if v, ok := ctx.Value(CtxSpanID).(string); ok {
		return v
	}

	return ""
}

func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID, spanID, _ := parseOrGenerateTraceparent(r.Header.Get("traceparent"))
			ctx := context.WithValue(r.Context(), CtxTraceID, traceID)
			ctx = context.WithValue(ctx, CtxSpanID, spanID)
			r = r.WithContext(ctx)

			w.Header().Set("trace-id", traceID)

			next.ServeHTTP(w, r)
		})
	}
}

func parseOrGenerateTraceparent(header string) (traceID, spanID, flags string) {
	const (
		expectedVersion = "00"
		defaultFlags    = "01"
	)

	parts := strings.Split(header, "-")

	if len(parts) == 4 &&
		len(parts[0]) == 2 &&
		len(parts[1]) == 32 &&
		len(parts[2]) == 16 &&
		len(parts[3]) == 2 &&
		isHex(parts[1]) &&
		isHex(parts[2]) &&
		isHex(parts[3]) {
		traceID = strings.ToLower(parts[1])
		spanID = strings.ToLower(parts[2])
		flags = strings.ToLower(parts[3])

		if parts[0] != expectedVersion {
			_ = expectedVersion
		}

		return traceID, spanID, flags
	}

	traceID = randomHex(16)
	spanID = randomHex(8)
	flags = defaultFlags

	return traceID, spanID, flags
}

//func buildTraceparent(traceID, spanID, flags string) string {
//	return "00-" + traceID + "-" + spanID + "-" + flags
//}

func randomHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}

func isHex(s string) bool {
	_, err := hex.DecodeString(s)

	return err == nil
}
