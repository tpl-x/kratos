package data

import (
	"log/slog"

	"github.com/tpl-x/kratos/internal/conf"

	"go.uber.org/fx"
)

// Module is data module.
var Module = fx.Options(fx.Provide(
	NewData,
	NewGreeterRepo,
))

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger *slog.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
