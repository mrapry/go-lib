package sdk

import (
	auth "github.com/mrapry/go-lib/sdk/auth-service"
	master_service "github.com/mrapry/go-lib/sdk/master-service"
	notification_service "github.com/mrapry/go-lib/sdk/notification-service"
)

// Option func type
type Option func(*sdkInstance)

// SetAuthService option func
func SetAuthService(authService auth.ServiceAuth) Option {
	return func(s *sdkInstance) {
		s.authService = authService
	}
}

// SetMasterService option func
func SetMasterService(masterService master_service.MasterService) Option {
	return func(s *sdkInstance) {
		s.masterService = masterService
	}
}

// SetNotificationService option func
func SetNotificationService(notificationService notification_service.ServiceNotification) Option {
	return func(s *sdkInstance) {
		s.notificationService = notificationService
	}
}

// SDK instance abstraction
type SDK interface {
	AuthService() auth.ServiceAuth
	MasterService() master_service.MasterService
	NotificationService() notification_service.ServiceNotification
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
	authService         auth.ServiceAuth
	masterService       master_service.MasterService
	notificationService notification_service.ServiceNotification
}

func (s *sdkInstance) AuthService() auth.ServiceAuth {
	return s.authService
}

func (s *sdkInstance) MasterService() master_service.MasterService {
	return s.masterService
}

func (s *sdkInstance) NotificationService() notification_service.ServiceNotification {
	return s.notificationService
}
