package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/victordinh/gokit-grpc/service"
)

type ValidateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ValidateUserResponse struct {
	Token string `json:"token"`
	Err   string `json:"err,omitempty"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Email string `json:"email"`
	Err   string `json:"err,omitempty"`
}

func makeValidateUserEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateUserRequest)
		token, err := svc.ValidateUser(ctx, req.Email, req.Password)
		if err != nil {
			return ValidateUserResponse{"", err.Error()}, err
		}

		return ValidateUserResponse{token, ""}, err
	}
}

func makeValidateTokenEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateTokenRequest)
		email, err := svc.ValidateToken(ctx, req.Token)
		if err != nil {
			return ValidateTokenResponse{"", err.Error()}, err
		}

		return ValidateTokenResponse{email, ""}, err
	}
}
