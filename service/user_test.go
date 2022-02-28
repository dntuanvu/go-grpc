package service

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func TestValidateUser(t *testing.T) {
	logger := log.NewLogfmtLogger(os.Stderr)
	service := NewService(logger)
	t.Run("invalid user", func(t *testing.T) {
		_, err := service.ValidateUser(context.Background(), "dntuanvu@gmail.com", "invalid")
		assert.NotNil(t, err)
		assert.Equal(t, "invalid_user", err.Error())
	})

	t.Run("invalid token", func(t *testing.T) {
		token, err := service.ValidateUser(context.Background(), "dntuanvu@gmail.com", "1234567")
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})
}
