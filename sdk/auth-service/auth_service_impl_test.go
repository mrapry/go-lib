package auth_service

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v5"

	"github.com/google/go-querystring/query"
	"github.com/jarcoal/httpmock"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/stretchr/testify/assert"
)

var (
	urlMock  = "http://auth.pln.tests"
	passMock = gofakeit.Password(true, true, true, true, false, 10)
)

func TestNewAuthService(t *testing.T) {
	testName := "Test new auth service"

	t.Run(testName, func(t *testing.T) {
		// set service
		service := NewAuthService(urlMock, passMock)

		assert.NotNil(t, service)
	})
}

func TestAuthServiceImplGenerateToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := map[string]struct {
		wantError bool
		wantMock  bool
		status    int
		response  interface{}
	}{
		"Test #1 positive auth service generate token": {
			wantError: false,
			wantMock:  true,
			response:  map[string]interface{}{"success": true, "message": "success", "code": http.StatusOK},
			status:    http.StatusOK,
		},
		"Test #2 negative auth service generate token": {
			wantError: true,
			wantMock:  false,
			response:  map[string]interface{}{"success": false, "message": golibshared.ErrorGeneral, "code": http.StatusBadGateway},
			status:    http.StatusBadGateway,
		},
		"Test #3 negative auth service generate token unmarshal data": {
			wantError: true,
			wantMock:  true,
			response:  map[string]interface{}{"message": false},
			status:    http.StatusBadGateway,
		},
	}

	// set service
	service := NewAuthService(urlMock, passMock)

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockUrl := fmt.Sprintf("%s/api/auth/create-token", urlMock)

			if test.wantMock {
				golibshared.CreateHttpRequestMock(http.MethodPost, mockUrl, test.status, test.response)
			} else {
				httpmock.Reset()
			}

			err := <-service.GenerateToken(context.Background(), GenerateTokenRequest{})
			if !test.wantError {
				assert.NoError(t, err.Error)
			} else {
				assert.Error(t, err.Error)
			}
		})
	}
}

func TestAuthServiceImplValidateToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := map[string]struct {
		wantError bool
		wantMock  bool
		status    int
		response  interface{}
	}{
		"Test #1 positive auth service validate token": {
			wantError: false,
			wantMock:  true,
			response:  map[string]interface{}{"success": true, "message": "success", "code": http.StatusOK},
			status:    http.StatusOK,
		},
		"Test #2 negative auth service validate token": {
			wantError: true,
			wantMock:  false,
			response:  map[string]interface{}{"success": false, "message": golibshared.ErrorGeneral, "code": http.StatusBadGateway},
			status:    http.StatusBadGateway,
		},
		"Test #3 negative auth service validate token unmarshal data": {
			wantError: true,
			wantMock:  true,
			response:  map[string]interface{}{"message": false},
			status:    http.StatusBadGateway,
		},
	}

	// set service
	service := NewAuthService(urlMock, passMock)

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockUrl := fmt.Sprintf("%s/api/auth/validate", urlMock)

			if test.wantMock {
				golibshared.CreateHttpRequestMock(http.MethodPost, mockUrl, test.status, test.response)
			} else {
				httpmock.Reset()
			}

			err := <-service.Validate(context.Background(), "")
			if !test.wantError {
				assert.NoError(t, err.Error)
			} else {
				assert.Error(t, err.Error)
			}
		})
	}
}

func TestAuthServiceImplRefreshToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := map[string]struct {
		wantError bool
		wantMock  bool
		status    int
		response  interface{}
	}{
		"Test #1 positive auth service refresh token": {
			wantError: false,
			wantMock:  true,
			response:  map[string]interface{}{"success": true, "message": "success", "code": http.StatusOK},
			status:    http.StatusOK,
		},
		"Test #2 negative auth service refresh token": {
			wantError: true,
			wantMock:  false,
			response:  map[string]interface{}{"success": false, "message": golibshared.ErrorGeneral, "code": http.StatusBadGateway},
			status:    http.StatusBadGateway,
		},
		"Test #3 negative auth service refresh token unmarshal data": {
			wantError: true,
			wantMock:  true,
			response:  map[string]interface{}{"message": false},
			status:    http.StatusBadGateway,
		},
	}

	// set service
	service := NewAuthService(urlMock, passMock)

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockUrl := fmt.Sprintf("%s/api/auth/refresh-token", urlMock)

			// create query
			param, _ := query.Values(TokenRequest{})
			mockUrl = fmt.Sprintf("%s?%s", mockUrl, param.Encode())

			if test.wantMock {
				golibshared.CreateHttpRequestMock(http.MethodGet, mockUrl, test.status, test.response)
			} else {
				httpmock.Reset()
			}

			err := <-service.RefreshToken(context.Background(), TokenRequest{})
			if !test.wantError {
				assert.NoError(t, err.Error)
			} else {
				assert.Error(t, err.Error)
			}
		})
	}
}

func TestAuthServiceImplRevokeToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := map[string]struct {
		wantError bool
		wantMock  bool
		status    int
		response  interface{}
	}{
		"Test #1 positive auth service revoke token": {
			wantError: false,
			wantMock:  true,
			response:  map[string]interface{}{"success": true, "message": "success", "code": http.StatusOK},
			status:    http.StatusOK,
		},
		"Test #2 negative auth service revoke token": {
			wantError: true,
			wantMock:  false,
			response:  map[string]interface{}{"success": false, "message": golibshared.ErrorGeneral, "code": http.StatusBadGateway},
			status:    http.StatusBadGateway,
		},
		"Test #3 negative auth service revoke token unmarshal data": {
			wantError: true,
			wantMock:  true,
			response:  map[string]interface{}{"message": false},
			status:    http.StatusBadGateway,
		},
	}

	// set service
	service := NewAuthService(urlMock, passMock)

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockUrl := fmt.Sprintf("%s/api/auth/delete-token", urlMock)

			// create query
			param, _ := query.Values(TokenRequest{})
			mockUrl = fmt.Sprintf("%s?%s", mockUrl, param.Encode())

			if test.wantMock {
				golibshared.CreateHttpRequestMock(http.MethodDelete, mockUrl, test.status, test.response)
			} else {
				httpmock.Reset()
			}

			err := <-service.RevokeToken(context.Background(), TokenRequest{})
			if !test.wantError {
				assert.NoError(t, err.Error)
			} else {
				assert.Error(t, err.Error)
			}
		})
	}
}
