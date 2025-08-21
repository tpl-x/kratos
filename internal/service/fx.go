package service

import "go.uber.org/fx"

// Module is service module.
var Module = fx.Options(fx.Provide(
	NewGreeterService,
))
