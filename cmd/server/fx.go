package main

import (
	"github.com/bufbuild/protovalidate-go"
	zapv2 "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tpl-x/kratos/internal/conf"
)

// Provider functions for configuration
func provideBootstrap() (*conf.Bootstrap, error) {
	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)

	if err := c.Load(); err != nil {
		return nil, err
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}

	return &bc, nil
}

// Provider function for validator
func provideValidator(bootstrap *conf.Bootstrap) (protovalidate.Validator, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	// Validate configuration
	if err := validator.Validate(bootstrap); err != nil {
		return nil, err
	}

	return validator, nil
}

// Provider functions for configuration sub-items
func provideDataConfig(bootstrap *conf.Bootstrap) *conf.Data {
	return bootstrap.Data
}

func provideLogConfig(bootstrap *conf.Bootstrap) *conf.Log {
	return bootstrap.Log
}

func provideServerConfig(bootstrap *conf.Bootstrap) *conf.Server {
	return bootstrap.Server
}

// Provider function for logger with service information
func provideLogger(zapLogger *zapv2.Logger) log.Logger {
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

// Provider function for Kratos application
func provideKratosApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	)
}

// Application start hook
func onStart(app *kratos.App) error {
	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()
	return nil
}

// Application stop hook
func onStop(app *kratos.App) error {
	return app.Stop()
}
