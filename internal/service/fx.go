package service

import "go.uber.org/fx"

// ProviderSet is service providers.
var ProviderSet = fx.Options(fx.Provide(
	NewGreeterService,
))
