package service

import (
	"context"
	"fmt"
	"github.com/bufbuild/protovalidate-go"

	v1 "github.com/tpl-x/kratos/api/helloworld/v1"
	"github.com/tpl-x/kratos/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServiceServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	v, err := protovalidate.New()
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}
	if err = v.Validate(in); err != nil {
		fmt.Println("validation failed:", err)
	} else {
		fmt.Println("validation succeeded")
	}
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.SayHelloResponse{Message: "Hello " + g.Hello}, nil
}
