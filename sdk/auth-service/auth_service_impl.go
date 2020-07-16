package auth_service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/mrapry/go-lib/golibutil"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/tracer"
	"github.com/mrapry/go-lib/wrapper"

	"github.com/google/go-querystring/query"
	"go.uber.org/zap/zapcore"
)

type authServiceImpl struct {
	host        string
	basic       string
	httpRequest golibutil.HTTPRequest
}

// NewAuthService constructor
func NewAuthService(host string, basic string) ServiceAuth {

	return &authServiceImpl{
		host:        host,
		basic:       basic,
		httpRequest: golibutil.NewHTTPRequest(5, 100*time.Millisecond),
	}
}
func (a *authServiceImpl) GenerateToken(ctx context.Context, payload GenerateTokenRequest) <-chan golibshared.Result {
	output := make(chan golibshared.Result)
	opName := "auth_service.generate_token"

	go tracer.WithTraceFunc(ctx, opName, func(c context.Context, tag map[string]interface{}) {
		defer close(output)

		var (
			response wrapper.HTTPResponse
			url      = fmt.Sprintf("%s/api/auth/create-token", a.host)
		)

		// set header
		header := map[string]string{
			echo.HeaderAuthorization: a.basic,
			echo.HeaderContentType:   echo.MIMEApplicationJSON,
		}

		// set payload
		req, _ := json.Marshal(payload)

		// http request
		resp, err := a.httpRequest.Do(ctx, http.MethodPost, url, req, header)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "http_request")
			output <- golibshared.Result{Error: err}
			return
		}

		// unmarshal to our target
		err = json.Unmarshal(resp, &response)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "unmarshal_auth_response")
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: response}
	})

	return output
}

func (a *authServiceImpl) Validate(ctx context.Context, token string) <-chan golibshared.Result {
	output := make(chan golibshared.Result)
	opName := "auth_service.validate_token"

	go tracer.WithTraceFunc(ctx, opName, func(c context.Context, tag map[string]interface{}) {
		defer close(output)

		var (
			response wrapper.HTTPResponse
			url      = fmt.Sprintf("%s/api/auth/validate", a.host)
		)

		// set header
		header := map[string]string{
			echo.HeaderAuthorization: a.basic,
			echo.HeaderContentType:   echo.MIMEApplicationJSON,
		}

		// set payload
		payload := ValidateTokenRequest{
			Token: token,
		}
		req, _ := json.Marshal(payload)

		// http request
		resp, err := a.httpRequest.Do(ctx, http.MethodPost, url, req, header)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "http_request")
			output <- golibshared.Result{Error: err}
			return
		}

		// unmarshal to our target
		err = json.Unmarshal(resp, &response)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "unmarshal_auth_response")
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: response}
	})

	return output
}

func (a *authServiceImpl) RefreshToken(ctx context.Context, payload TokenRequest) <-chan golibshared.Result {
	output := make(chan golibshared.Result)
	opName := "auth_service.refresh_token"

	go tracer.WithTraceFunc(ctx, opName, func(c context.Context, tag map[string]interface{}) {
		defer close(output)

		var (
			response wrapper.HTTPResponse
			url      = fmt.Sprintf("%s/api/auth/refresh-token", a.host)
		)

		// set header
		header := map[string]string{
			echo.HeaderAuthorization: a.basic,
		}

		// create query
		param, _ := query.Values(payload)
		url = fmt.Sprintf("%s?%s", url, param.Encode())

		// http request
		resp, err := a.httpRequest.Do(ctx, http.MethodGet, url, nil, header)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "http_request")
			output <- golibshared.Result{Error: err}
			return
		}

		// unmarshal to our target
		err = json.Unmarshal(resp, &response)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "unmarshal_auth_response")
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: response}
	})

	return output
}

func (a *authServiceImpl) RevokeToken(ctx context.Context, payload TokenRequest) <-chan golibshared.Result {
	output := make(chan golibshared.Result)
	opName := "auth_service.revoke_token"

	go tracer.WithTraceFunc(ctx, opName, func(c context.Context, tag map[string]interface{}) {
		defer close(output)

		var (
			response wrapper.HTTPResponse
			url      = fmt.Sprintf("%s/api/auth/delete-token", a.host)
		)

		// set header
		header := map[string]string{
			echo.HeaderAuthorization: a.basic,
		}

		// create query
		param, _ := query.Values(payload)
		url = fmt.Sprintf("%s?%s", url, param.Encode())

		// http request
		resp, err := a.httpRequest.Do(ctx, http.MethodDelete, url, nil, header)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "http_request")
			output <- golibshared.Result{Error: err}
			return
		}

		// unmarshal to our target
		err = json.Unmarshal(resp, &response)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "unmarshal_auth_response")
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: response}
	})

	return output
}
