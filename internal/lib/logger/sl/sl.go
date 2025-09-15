package sl

import (
	"context"
	"log/slog"
	"plata_currency_quotation/internal/lib/http-server/middleware/traceparent"
)

var Log *slog.Logger

func Err(err error) slog.Attr {
	return slog.Attr{Key: "error", Value: slog.StringValue(err.Error())}
}

func TraceId(ctx context.Context) slog.Attr {
	return slog.Attr{Key: string(traceparent.CtxTraceID), Value: slog.StringValue(ctx.Value(traceparent.CtxTraceID).(string))}
}
