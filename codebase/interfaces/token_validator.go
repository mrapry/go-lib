package interfaces

import (
	"context"

	"github.com/mrapry/go-lib/golibshared"
)

// TokenValidator abstract interface
type TokenValidator interface {
	Validate(ctx context.Context, token string) <-chan golibshared.Result
}
