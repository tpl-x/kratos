package biz

import (
	"context"
	"log/slog"

	v1 "github.com/tpl-x/kratos/api/helloworld/v1"

	"github.com/go-kratos/kratos/v3/errors"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Greeter struct {
	Hello string
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Save(context.Context, *Greeter) (*Greeter, error)
	Update(context.Context, *Greeter) (*Greeter, error)
	FindByID(context.Context, int64) (*Greeter, error)
	ListByHello(context.Context, string) ([]*Greeter, error)
	ListAll(context.Context) ([]*Greeter, error)
}

// GreeterUseCase is a Greeter useCase.
type GreeterUseCase struct {
	repo GreeterRepo
	log  *slog.Logger
}

// NewGreeterUseCase new a Greeter useCase.
func NewGreeterUseCase(repo GreeterRepo, logger *slog.Logger) *GreeterUseCase {
	return &GreeterUseCase{
		repo: repo,
		log:  logger.With(slog.String("module", "biz/greeterUseCase")),
	}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterUseCase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.InfoContext(ctx, "CreateGreeter", slog.String("hello", g.Hello))
	return uc.repo.Save(ctx, g)
}
