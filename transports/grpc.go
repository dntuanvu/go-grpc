package transports

import (
	"context"

	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"

	"github.com/victordinh/gokit-grpc/endpoints"
	"github.com/victordinh/gokit-grpc/pb"
)

type gRPCServer struct {
	add           gt.Handler
	validateUser  gt.Handler
	validateToken gt.Handler
}

func NewGRPCServer(endpoints endpoints.Endpoints, logger log.Logger) pb.VictorServiceServer {
	return &gRPCServer{
		add: gt.NewServer(
			endpoints.Add,
			decodeMathRequest,
			encodeMathResponse,
		),
		validateUser: gt.NewServer(
			endpoints.ValidateUserEndpoint,
			decodeValidateUserRequestGRPC,
			encodeValidateUserResponse,
		),
		validateToken: gt.NewServer(
			endpoints.ValidateTokenEndpoint,
			decodeValidateTokenRequestGRPC,
			encodeValidateTokenResponse,
		),
	}
}

func (s *gRPCServer) Add(ctx context.Context, req *pb.MathRequest) (*pb.MathResponse, error) {
	_, resp, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.MathResponse), nil
}

func (s *gRPCServer) ValidateUser(ctx context.Context, req *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	_, resp, err := s.validateUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ValidateUserResponse), nil
}

func (s *gRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	_, resp, err := s.validateToken.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ValidateTokenResponse), nil
}

func decodeMathRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.MathRequest)
	return endpoints.MathReq{NumA: req.NumA, NumB: req.NumB}, nil
}

func encodeMathResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.MathResp)
	return &pb.MathResponse{Result: resp.Result}, nil
}

func decodeValidateUserRequestGRPC(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ValidateUserRequest)
	return endpoints.ValidateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func encodeValidateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.ValidateUserResponse)
	return &pb.ValidateUserResponse{
		Token: resp.Token,
	}, nil
}

func decodeValidateTokenRequestGRPC(_ context.Context, response interface{}) (interface{}, error) {
	req := response.(*pb.ValidateTokenRequest)
	return &pb.ValidateTokenRequest{
		Token: req.Token,
	}, nil
}

func encodeValidateTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.ValidateTokenResponse)
	return &pb.ValidateTokenResponse{
		Email: resp.Email,
	}, nil
}
