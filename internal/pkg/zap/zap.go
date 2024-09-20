package zap

import (
	zapmiddleware "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/tpl-x/kratos/api/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// NewLoggerWithLumberjack returns a HandlerFunc that adds a zap logger to the context.
func NewLoggerWithLumberjack(bc *conf.Bootstrap) *zapmiddleware.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	lumberjackLogger := &lumberjack.Logger{
		Filename:   bc.Log.LogPath,
		MaxSize:    int(bc.Log.MaxSize),
		MaxBackups: int(bc.Log.MaxKeepFiles),
		MaxAge:     int(bc.Log.MaxKeepDays),
		Compress:   bc.Log.Compress,
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberjackLogger),
	)
	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	var logLevel zapcore.Level
	logLevel = convertInnerLogLevelToZapLogLevel(bc.Log.LogLevel, logLevel)
	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		zap.NewAtomicLevelAt(logLevel),
	)
	zapCore := zap.New(core, zap.AddCaller())
	zapLog := zapmiddleware.NewLogger(zapCore)
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
