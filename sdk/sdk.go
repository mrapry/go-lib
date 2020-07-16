package sdk

import (
	auth "github.com/mrapry/go-lib/sdk/auth-service"
	caterpillar "github.com/mrapry/go-lib/sdk/caterpillar-service"
	dil "github.com/mrapry/go-lib/sdk/dil-service"
)

// Option func type
type Option func(*SDK)

// SetAuthService option func
func SetAuthService(authService *auth.ServiceAuth) Option {
	return func(s *SDK) {
		s.AuthService = *authService
	}
}

// SetDilService option func
func SetDilService(dilService dil.ServiceDil) Option {
	return func(s *SDK) {
		s.DilService = dilService
	}
}

// SetCaterpillarService option func
func SetCaterpillarService(caterpillarService caterpillar.ServiceCaterpillar) Option {
	return func(s *SDK) {
		s.CaterpillarService = caterpillarService
	}
}

// SDK instance
type SDK struct {
	AuthService        auth.ServiceAuth
	DilService         dil.ServiceDil
	CaterpillarService caterpillar.ServiceCaterpillar
}

// NewSDK constructor with each service option
func NewSDK(opts ...Option) *SDK {
	sdk := new(SDK)

	for _, o := range opts {
		o(sdk)
	}

	return sdk
}
