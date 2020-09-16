package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/labstack/echo"
	"github.com/mrapry/go-lib/golibshared"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {

	midd := &Middleware{
		username: "user", password: "da1c25d8-37c8-41b1-afe2-42dd4825bfea",
	}

	t.Run("Test With Valid Auth", func(t *testing.T) {

		err := midd.Basic(context.Background(), "dXNlcjpkYTFjMjVkOC0zN2M4LTQxYjEtYWZlMi00MmRkNDgyNWJmZWE=")
		assert.NoError(t, err)
	})

	t.Run("Test With Invalid Auth #1", func(t *testing.T) {

		err := midd.Basic(context.Background(), "MjIyMjphc2RzZA==")
		assert.Error(t, err)
	})

	t.Run("Test With Invalid Auth #2", func(t *testing.T) {

		err := midd.Basic(context.Background(), "Basic")
		assert.Error(t, err)
	})

	t.Run("Test With Invalid Auth #3", func(t *testing.T) {

		err := midd.Basic(context.Background(), "Bearer xxx")
		assert.Error(t, err)
	})

	t.Run("Test With Invalid Auth #4", func(t *testing.T) {

		err := midd.Basic(context.Background(), "zzzzzzz")
		assert.Error(t, err)
	})

	t.Run("Test With Invalid Auth #5", func(t *testing.T) {

		err := midd.Basic(context.Background(), "Basic dGVzdGluZw==")
		assert.Error(t, err)
	})
}

func TestMiddlewareHTTPBasicAuth(t *testing.T) {
	testCase := map[string]struct {
		authKey   string
		authValue string
	}{
		"Test #1 positive http basic auth": {
			authKey:   echo.HeaderAuthorization,
			authValue: fmt.Sprintf("%s %s", strings.ToTitle(Basic), "dXNlcjpkYTFjMjVkOC0zN2M4LTQxYjEtYWZlMi00MmRkNDgyNWJmZWE="),
		},
		"Test #2 negative false http basic auth": {
			authKey:   echo.HeaderAuthorization,
			authValue: fmt.Sprintf("%s %s", strings.ToTitle(Basic), gofakeit.Word()),
		},
		"Test #3 negative no basic http basic auth": {
			authKey:   echo.HeaderAuthorization,
			authValue: fmt.Sprintf("%s", strings.ToTitle(Basic)),
		},
		"Test #4 negative no authorization http basic auth": {},
	}

	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			// set middleware
			mw := &Middleware{
				username: "user", password: "da1c25d8-37c8-41b1-afe2-42dd4825bfea",
			}

			// set headers
			headers := map[string]string{
				test.authKey: test.authValue,
			}

			// mock http echo
			c := golibshared.SetEchoHTTPMock(`/`, http.MethodGet, ``, headers)

			// set handler
			h := mw.HTTPBasicAuth(true)(func(c echo.Context) error {
				return c.String(http.StatusOK, gofakeit.Word())
			})

			h(c)
		})
	}
}
