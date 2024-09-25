package zap

import (
	zapv2 "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/tpl-x/kratos/internal/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// NewLoggerWithLumberjack returns a HandlerFunc that adds a zap logger to the context.
func NewLoggerWithLumberjack(logConfig *conf.Log) *zapv2.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logConfig.LogPath,
		MaxSize:    int(logConfig.MaxSize),
		MaxBackups: int(logConfig.MaxKeepFiles),
		MaxAge:     int(logConfig.MaxKeepDays),
		Compress:   logConfig.Compress,
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberjackLogger),
	)
	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	var logLevel zapcore.Level
	logLevel = convertInnerLogLevelToZapLogLevel(logConfig.LogLevel, logLevel)
	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		zap.NewAtomicLevelAt(logLevel),
	)
	zapCore := zap.New(core, zap.AddCaller())
	zapLog := zapv2.NewLogger(zapCore)
	defer func() { _ = zapLog.Sync() }()
	return zapLog
}

// convertServerLogLevel converts server log level to zap log level.
func convertInnerLogLevelToZapLogLevel(svrLogLevel conf.LogLevel, logLevel zapcore.Level) zapcore.Level {
	switch svrLogLevel {
	case conf.LogLevel_Debug:
		logLevel = zapcore.DebugLevel
	case conf.LogLevel_Info:
		logLevel = zapcore.InfoLevel
	case conf.LogLevel_Warn:
		logLevel = zapcore.WarnLevel
	case conf.LogLevel_Error:
		logLevel = zapcore.ErrorLevel
	case conf.LogLevel_Fatal:
		logLevel = zapcore.FatalLevel
	}
	return logLevel
}
