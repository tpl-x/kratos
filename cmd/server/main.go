package main

import (
	"context"
	"flag"
	"github.com/tpl-x/kratos/internal/biz"
	"github.com/tpl-x/kratos/internal/data"
	"github.com/tpl-x/kratos/internal/pkg/zap"
	"github.com/tpl-x/kratos/internal/server"
	"github.com/tpl-x/kratos/internal/service"
	"go.uber.org/fx"
	"os"

	"github.com/go-kratos/kratos/v2"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagConf is the config flag.
	flagConf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	// Create fx application
	app := fx.New(
		// Provide basic dependencies
		fx.Provide(
			provideBootstrap,
			provideValidator,
			provideDataConfig,
			provideLogConfig,
			provideServerConfig,
		),

		// Provide logging related dependencies
		fx.Provide(
			zap.NewLoggerWithLumberjack,
			provideLogger,
		),

		// Include ProviderSets from other modules
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,

		// Provide Kratos application
		fx.Provide(provideKratosApp),

		// Set up lifecycle hooks
		fx.Invoke(func(lc fx.Lifecycle, app *kratos.App) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					return onStart(app)
				},
				OnStop: func(context.Context) error {
					return onStop(app)
				},
			})
		}),
	)

	// Run the application
	app.Run()
}
