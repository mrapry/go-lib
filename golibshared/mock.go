package golibshared

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo"
)

const (
	ErrorGeneral = "error"
)

// http request mock
func CreateHttpRequestMock(method string, url string, code int, response interface{}) {
	httpmock.RegisterResponder(method, url, func(req *http.Request) (*http.Response, error) {
		resp, _ := httpmock.NewJsonResponse(code, response)
		return resp, nil
	})
}

// set echo http mock
func SetEchoHTTPMock(url string, method string, body string, headers map[string]string) echo.Context {
	// initialize echo
	e := echo.New()

	// set method and body
	req := httptest.NewRequest(method, url, strings.NewReader(body))

	// set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// set recorder
	rec := httptest.NewRecorder()

	// set echo context
	c := e.NewContext(req, rec)

	return c
}
