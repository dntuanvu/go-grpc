package service

import (
	"context"

	"github.com/victordinh/gokit-grpc/auth/security"
)

func (s service) ValidateUser(ctx context.Context, email, password string) (string, error) {
	if email == "dntuanvu@gmail.com" && password != "1234567" {
		return "nil", ErrInvalidUser
	}

	token, err := security.NewToken(email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s service) ValidateToken(ctx context.Context, token string) (string, error) {
	t, err := security.ParseToken(token)
	if err != nil {
		return "", ErrInvalidToken
	}

	tData, err := security.GetClaims(t)
	if err != nil {
		return "", ErrInvalidToken
	}

	return tData["email"].(string), nil
}
