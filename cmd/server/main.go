package main

import (
	"flag"
	"github.com/tpl-x/kratos/internal/biz"
	"github.com/tpl-x/kratos/internal/data"
	"github.com/tpl-x/kratos/internal/server"
	"github.com/tpl-x/kratos/internal/service"
	"go.uber.org/fx"
	"os"

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
		// Provide configs
		fx.Provide(
			provideConfigs,
		),
		// Provide logging related dependencies
		loggingModule,

		// Include other modules
		server.Module,
		data.Module,
		biz.Module,
		service.Module,
		// Provide Kratos application
		appModule,
	)

	// Run the application
	app.Run()
}
