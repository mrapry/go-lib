package auth_service

import (
	"context"

	"github.com/mrapry/go-lib/golibshared"
)

// ServiceAuth interface
type ServiceAuth interface {
	GenerateToken(ctx context.Context, payload GenerateTokenRequest) <-chan golibshared.Result
	Validate(ctx context.Context, token string) <-chan golibshared.Result
	RefreshToken(ctx context.Context, payload TokenRequest) <-chan golibshared.Result
	RevokeToken(ctx context.Context, payload TokenRequest) <-chan golibshared.Result
}
