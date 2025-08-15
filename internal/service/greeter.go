package service

import (
	"context"
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/tpl-x/kratos/api/helloworld/v1"
	"github.com/tpl-x/kratos/internal/biz"
)

var _ v1.GreeterServiceServer = (*GreeterService)(nil)

// GreeterService is a greeter service.
type GreeterService struct {
	log       *log.Helper
	uc        *biz.GreeterUseCase
	validator protovalidate.Validator
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUseCase, logger log.Logger, validator protovalidate.Validator) *GreeterService {
	return &GreeterService{
		uc:        uc,
		validator: validator,
		log:       log.NewHelper(log.With(logger, "module", "service/greeterService")),
	}
}

// LuckySearch implements helloworld.LuckySearch
func (s *GreeterService) LuckySearch(ctx context.Context, request *v1.LuckySearchRequest) (*v1.LuckySearchResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		s.log.Error("request validate failed", err)
		return nil, err
	}
	s.log.Info("validation succeeded")

	keyword := request.GetKeyword()
	resp := &v1.LuckySearchResponse{
		RedirectTo: fmt.Sprintf("https://www.google.com/search?q=%s", keyword),
		StatusCode: 302,
	}
	return resp, nil
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, request *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		s.log.Error("request validate failed", err)
		return nil, err
	}
	s.log.Info("validation succeeded")
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: request.Name})
	if err != nil {
		return nil, err
	}
	return &v1.SayHelloResponse{Message: "Hello " + g.Hello}, nil
}
