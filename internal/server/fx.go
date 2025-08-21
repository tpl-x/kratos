package server

import (
	"go.uber.org/fx"
)

// Module is server module.
var Module = fx.Options(fx.Provide(
	NewGRPCServer,
	NewHTTPServer,
))
