package o11y

import (
	"context"
	"strconv"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func TraceFields(ctx context.Context) logrus.Fields {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.IsValid() {
		return logrus.Fields{
			"dd.trace_id": convertTraceID(spanCtx.TraceID().String()),
			"dd.span_id":  convertTraceID(spanCtx.SpanID().String()),
		}
	}
	return logrus.Fields{}
}

func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
