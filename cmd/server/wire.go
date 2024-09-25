//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/bufbuild/protovalidate-go"
	zapv2 "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/google/wire"
	"github.com/tpl-x/kratos/internal/biz"
	"github.com/tpl-x/kratos/internal/conf"
	"github.com/tpl-x/kratos/internal/data"
	"github.com/tpl-x/kratos/internal/pkg/zap"
	"github.com/tpl-x/kratos/internal/server"
	"github.com/tpl-x/kratos/internal/service"
)

func newZapLoggerWith(zapLogger *zapv2.Logger) log.Logger {
	return log.With(zapLogger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
}

var appSet = wire.NewSet(
	wire.FieldsOf(new(*conf.Bootstrap),
		"Data",
		"Log",
		"Server",
	),
	zap.NewLoggerWithLumberjack,
	newZapLoggerWith,
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, *protovalidate.Validator) (*kratos.App, func(), error) {
	panic(wire.Build(appSet, server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
