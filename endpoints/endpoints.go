package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/victordinh/gokit-grpc/service"
)

type Endpoints struct {
	HealthCheckEndpoint   endpoint.Endpoint
	Add                   endpoint.Endpoint
	ValidateUserEndpoint  endpoint.Endpoint
	ValidateTokenEndpoint endpoint.Endpoint
}

type MathReq struct {
	NumA float32 `json:"numA"`
	NumB float32 `json:"numB"`
}

type MathResp struct {
	Result float32
}

type HealthCheckResp struct {
	Status string `json:"status"`
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		HealthCheckEndpoint:   makeHealthCheckEndpoint(s),
		Add:                   makeAddEndpoints(s),
		ValidateUserEndpoint:  makeValidateUserEndpoint(s),
		ValidateTokenEndpoint: makeValidateTokenEndpoint(s),
	}
}

func makeHealthCheckEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthCheckResp{Status: "ok"}, nil
	}
}

func makeAddEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathReq)
		result, _ := s.Add(ctx, req.NumA, req.NumB)
		return MathResp{Result: result}, nil
	}
}
