package service

import (
	"context"
	"errors"

	"github.com/go-kit/log"
)

type service struct {
	logger log.Logger
}

var (
	ErrInvalidUser  = errors.New("invalid_user")
	ErrInvalidToken = errors.New("invalid_token")
)

type Service interface {
	Add(ctx context.Context, numA, numB float32) (float32, error)

	ValidateUser(ctx context.Context, email, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

func NewService(logger log.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (s service) Add(ctx context.Context, numA, numB float32) (float32, error) {
	return numA + numB, nil
}
