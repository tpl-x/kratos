package server

import (
	"go.uber.org/fx"
)

// ProviderSet is server providers.
var ProviderSet = fx.Options(fx.Provide(
	NewGRPCServer,
	NewHTTPServer,
))
