package endpoints

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/victordinh/gokit-grpc/service"
)

func TestMakeValidateUserEndpoint(t *testing.T) {
	logger := log.NewLogfmtLogger(os.Stderr)
	s := service.NewService(logger)
	endpoint := makeValidateUserEndpoint(s)
	t.Run("valid user", func(t *testing.T) {
		req := ValidateUserRequest{
			Email:    "dntuanvu@gmail.com",
			Password: "1234567",
		}

		_, err := endpoint(context.Background(), req)
		if err != nil {
			t.Errorf("expected %v received %v", nil, err)
		}
	})

	t.Run("invalid user", func(t *testing.T) {
		req := ValidateUserRequest{
			Email:    "dntuanvu@gmail.com",
			Password: "123456",
		}

		_, err := endpoint(context.Background(), req)
		if err == nil {
			t.Errorf("expected %v received %v", service.ErrInvalidUser, err)
		}
	})

}
