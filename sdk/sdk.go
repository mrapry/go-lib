package sdk

import (
	auth "github.com/mrapry/go-lib/sdk/auth-service"
)

// Option func type
type Option func(*sdkInstance)

// SetAuthService option func
func SetAuthService(authService auth.ServiceAuth) Option {
	return func(s *sdkInstance) {
		s.authService = authService
	}
}

// SDK instance abstraction
type SDK interface {
	AuthService() auth.ServiceAuth
}

// NewSDK constructor with each service option
func NewSDK(opts ...Option) SDK {
	sdk := new(sdkInstance)

	for _, o := range opts {
		o(sdk)
	}

	return sdk
}

type sdkInstance struct {
	authService auth.ServiceAuth
}

func (s *sdkInstance) AuthService() auth.ServiceAuth {
	return s.authService
}
