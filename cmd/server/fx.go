package main

import (
	"buf.build/go/protovalidate"
	"context"
	zapv2 "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tpl-x/kratos/internal/conf"
	"github.com/tpl-x/kratos/internal/pkg/zap"
	"go.uber.org/fx"
)

type ConfigBundle struct {
	fx.Out

	Bootstrap *conf.Bootstrap
	Data      *conf.Data
	Log       *conf.Log
	Server    *conf.Server
	Validator protovalidate.Validator
}

func provideConfigs() (ConfigBundle, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return ConfigBundle{}, err
	}
	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)

	if err := c.Load(); err != nil {
		return ConfigBundle{}, err
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return ConfigBundle{}, err
	}

	if err := validator.Validate(&bc); err != nil {
		return ConfigBundle{}, err
	}

	return ConfigBundle{
		Bootstrap: &bc,
		Data:      bc.Data,
		Log:       bc.Log,
		Server:    bc.Server,
		Validator: validator,
	}, nil
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

// newKratosApp function for Kratos application
func newKratosApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	)
}

func setupLifecycle(lc fx.Lifecycle, app *kratos.App) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return onStart(app)
		},
		OnStop: func(context.Context) error {
			return onStop(app)
		},
	})
}

var appModule = fx.Options(
	fx.Provide(newKratosApp),
	fx.Invoke(setupLifecycle),
)

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

var loggingModule = fx.Options(
	fx.Provide(
		zap.NewLoggerWithLumberjack,
		provideLogger,
	),
)
