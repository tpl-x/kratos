package biz

import "go.uber.org/fx"

// Module is biz module.
var Module = fx.Options(fx.Provide(
	NewGreeterUseCase,
))
