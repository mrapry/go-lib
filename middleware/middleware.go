package middleware

import (
	"context"

	"github.com/mrapry/go-lib/codebase/interfaces"
	"github.com/mrapry/go-lib/config"
	"github.com/mrapry/go-lib/golibshared"
)

// Middleware impl
type Middleware struct {
	tokenValidator      interfaces.TokenValidator
	username, password  string
	authTypeCheckerFunc map[string]func(context.Context, string) (*golibshared.TokenClaim, error)
}

// NewMiddleware create new middleware instance
func NewMiddleware(tokenValidator interfaces.TokenValidator) *Middleware {
	mw := &Middleware{
		tokenValidator: tokenValidator,
		username:       config.BaseEnv().BasicAuthUsername,
		password:       config.BaseEnv().BasicAuthPassword,
	}

	mw.authTypeCheckerFunc = map[string]func(context.Context, string) (*golibshared.TokenClaim, error){
		Basic: func(ctx context.Context, key string) (*golibshared.TokenClaim, error) {
			return nil, mw.Basic(ctx, key)
		},
		Bearer: func(ctx context.Context, token string) (*golibshared.TokenClaim, error) {
			return mw.Bearer(ctx, token)
		},
	}

	return mw
}
