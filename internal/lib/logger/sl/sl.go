package sl

import (
	"context"
	"log/slog"
	"plata_currency_quotation/internal/lib/http-server/middleware/trace-id"
)

var Log *slog.Logger

func Err(err error) slog.Attr {
	return slog.Attr{Key: "error", Value: slog.StringValue(err.Error())}
}

func TraceId(ctx context.Context) slog.Attr {
	return slog.Attr{Key: trace_id.CtxTraceId, Value: slog.StringValue(ctx.Value(trace_id.CtxTraceId).(string))}
}
