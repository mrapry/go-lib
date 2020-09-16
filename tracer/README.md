# How to use tracing using jaeger tracing
example in your `main.go`:
```golang
package main

import (
    // import go-lib from mrapry go-lib
    "github.com/mrapry/go-lib/tracer"
)

func main() {
    serviceName := strings.TrimSuffix("{servicename}-"+os.Getenv("SERVER_ENV"), "{-production or -dev or -stg}")
    tracer.InitOpenTracing(os.Getenv("JAEGER_HOST"), serviceName)
}
```

example in your `main_http.go`:
```golang
package main

import (
    // import go-lib from mrapry go-lib
    "github.com/mrapry/go-lib/tracer"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

func main_http() {
    // using echo
    e := echo.New()
    e.Use(echo.WrapMiddleware(tracer.Middleware))
    e.Use(middleware.BodyDump(func(c echo.Context, req []byte, res []byte) {
        span := opentracing.SpanFromContext(c.Request().Context())
        statusCode := c.Response().Status
        ext.HTTPStatusCode.Set(span, uint16(statusCode))
        if statusCode >= http.StatusBadRequest {
            ext.Error.Set(span, true)
        }
        span.SetTag("response.body", string(res))
    }))
}
```