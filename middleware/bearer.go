package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/mrapry/go-lib/golibhelper"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/mrapry/go-lib/wrapper"

	"github.com/labstack/echo"
)

const (
	Bearer = "bearer"
)

// Bearer token validator
func (m *Middleware) Bearer(ctx context.Context, tokenString string) (*golibshared.TokenClaim, error) {
	resp := <-m.tokenValidator.Validate(ctx, tokenString)

	if resp.Error != nil {
		return nil, resp.Error
	}

	tokenClaim, ok := resp.Data.(*golibshared.TokenClaim)
	if !ok {
		return nil, errors.New("Validate token: result is not claim data")
	}

	return tokenClaim, nil
}

// HTTPBearerAuth http jwt token middleware
func (m *Middleware) HTTPBearerAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authorization := c.Request().Header.Get(echo.HeaderAuthorization)
			if authorization == "" {
				return wrapper.NewHTTPResponse(http.StatusUnauthorized, "Invalid authorization").JSON(c.Response())
			}

			authValues := strings.Split(authorization, " ")
			authType := strings.ToLower(authValues[0])
			if authType != Bearer || len(authValues) != 2 {
				return wrapper.NewHTTPResponse(http.StatusUnauthorized, "Invalid authorization").JSON(c.Response())
			}

			tokenString := authValues[1]
			tokenClaim, err := m.Bearer(c.Request().Context(), tokenString)
			if err != nil {
				return wrapper.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c.Response())
			}

			c.Set(golibhelper.TokenClaimKey, tokenClaim)
			return next(c)
		}
	}
}

func (m *Middleware) GraphQLBearerAuth(ctx context.Context) *golibshared.TokenClaim {
	headers := ctx.Value(golibshared.ContextKey("headers")).(http.Header)
	authorization := headers.Get(echo.HeaderAuthorization)

	authValues := strings.Split(authorization, " ")
	authType := strings.ToLower(authValues[0])
	if authType != Bearer || len(authValues) != 2 {
		panic("Invalid authorization")
	}

	tokenClaim, err := m.Bearer(ctx, authValues[1])
	if err != nil {
		panic(err)
	}

	return tokenClaim
}
