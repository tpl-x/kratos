package logging

import (
	"context"
	"io"
	"log/slog"
	"os"

	klog "github.com/go-kratos/kratos/v3/log"
	"github.com/tpl-x/kratos/internal/conf"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLoggerWithLumberjack builds the process logger using Kratos v3 slog APIs.
func NewLoggerWithLumberjack(logConfig *conf.Log) *slog.Logger {
	writer := io.MultiWriter(
		os.Stdout,
		&lumberjack.Logger{
			Filename:   logConfig.LogPath,
			MaxSize:    int(logConfig.MaxSize),
			MaxBackups: int(logConfig.MaxKeepFiles),
			MaxAge:     int(logConfig.MaxKeepDays),
			Compress:   logConfig.Compress,
		},
	)

	return klog.NewLogger(
		slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level:     convertInnerLogLevelToSlogLevel(logConfig.LogLevel),
			AddSource: true,
		}),
		klog.WithExtractor(traceAttrs),
	)
}

func traceAttrs(ctx context.Context) []slog.Attr {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return nil
	}
	return []slog.Attr{
		slog.String("trace.id", spanContext.TraceID().String()),
		slog.String("span.id", spanContext.SpanID().String()),
	}
}

func convertInnerLogLevelToSlogLevel(svrLogLevel conf.LogLevel) slog.Level {
	switch svrLogLevel {
	case conf.LogLevel_Debug:
		return slog.LevelDebug
	case conf.LogLevel_Info:
		return slog.LevelInfo
	case conf.LogLevel_Warn:
		return slog.LevelWarn
	case conf.LogLevel_Error:
		return slog.LevelError
	case conf.LogLevel_Fatal:
		return klog.LevelFatal
	default:
		return slog.LevelInfo
	}
}
