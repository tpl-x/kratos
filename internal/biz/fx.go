package biz

import "go.uber.org/fx"

// ProviderSet is biz providers.
var ProviderSet = fx.Options(fx.Provide(
	NewGreeterUseCase,
))
