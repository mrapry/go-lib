package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mrapry/go-lib/golibhelper"
)

// Logger function for writing all request log into console
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		req := c.Request()
		res := c.Response()

		err := next(c)

		statusCode := res.Status
		if he, ok := err.(*echo.HTTPError); ok {
			statusCode = he.Code
		}
		end := time.Now()

		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(req.Method)
		resetColor := golibhelper.Reset

		fmt.Fprintf(os.Stdout, "%s[SERVICE-REST]%s :%s %v | %s %3d %s | %13v | %15s | %s %-7s %s %s\n",
			golibhelper.White, golibhelper.Reset, req.URL.Port(),
			end.Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, resetColor,
			end.Sub(start),
			c.RealIP(),
			methodColor, req.Method, resetColor,
			req.RequestURI,
		)
		return err
	}
}

func colorForStatus(code int) []byte {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return golibhelper.Green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return golibhelper.White
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return golibhelper.Yellow
	default:
		return golibhelper.Red
	}
}

func colorForMethod(method string) []byte {
	switch method {
	case "GET":
		return golibhelper.Blue
	case "POST":
		return golibhelper.Cyan
	case "PUT":
		return golibhelper.Yellow
	case "DELETE":
		return golibhelper.Red
	case "PATCH":
		return golibhelper.Green
	case "HEAD":
		return golibhelper.Magenta
	case "OPTIONS":
		return golibhelper.White
	default:
		return golibhelper.Reset
	}
}
