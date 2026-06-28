package main

import (
	"buf.build/go/protovalidate"
	"context"
	"log/slog"

	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/config"
	"github.com/go-kratos/kratos/v3/config/file"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/go-kratos/kratos/v3/transport/http"
	"github.com/tpl-x/kratos/internal/conf"
	applogging "github.com/tpl-x/kratos/internal/pkg/logging"
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

	serverCfg, err := config.Get[conf.Server](c, "server")
	if err != nil {
		return ConfigBundle{}, err
	}
	dataCfg, err := config.Get[conf.Data](c, "data")
	if err != nil {
		return ConfigBundle{}, err
	}
	logCfg, err := config.Get[conf.Log](c, "log")
	if err != nil {
		return ConfigBundle{}, err
	}

	bc := &conf.Bootstrap{
		Server: &serverCfg,
		Data:   &dataCfg,
		Log:    &logCfg,
	}
	if err := validator.Validate(bc); err != nil {
		return ConfigBundle{}, err
	}

	return ConfigBundle{
		Bootstrap: bc,
		Data:      bc.Data,
		Log:       bc.Log,
		Server:    bc.Server,
		Validator: validator,
	}, nil
}

// Provider function for logger with service information
func provideLogger(logConfig *conf.Log) *slog.Logger {
	logger := applogging.NewLoggerWithLumberjack(logConfig).With(
		slog.String("service.id", id),
		slog.String("service.name", Name),
		slog.String("service.version", Version),
	)
	log.SetDefault(logger)
	return logger
}

// newKratosApp function for Kratos application
func newKratosApp(logger *slog.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
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
		provideLogger,
	),
)
