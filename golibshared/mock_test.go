package golibshared

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestSetEchoHTTPMock(t *testing.T) {
	t.Run("Test set echo http mock", func(t *testing.T) {
		c := SetEchoHTTPMock("/", http.MethodGet, `{"message" : "success"}`, map[string]string{echo.HeaderContentType: echo.MIMEApplicationJSON})

		assert.NotNil(t, c)
	})
}

func TestCreateHttpRequestMock(t *testing.T) {
	urlMock := "http://mrapry.test.com"
	responseMock := map[string]string{"success": "true"}

	testCases := map[string]struct {
		url      string
		method   string
		code     int
		response map[string]string
	}{
		"Test #1 create http request mock get":    {urlMock, http.MethodGet, http.StatusOK, responseMock},
		"Test #2 create http request mock post":   {urlMock, http.MethodPost, http.StatusCreated, responseMock},
		"Test #3 create http request mock put":    {urlMock, http.MethodPut, http.StatusBadRequest, responseMock},
		"Test #4 create http request mock delete": {urlMock, http.MethodDelete, http.StatusOK, responseMock},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			CreateHttpRequestMock(test.url, test.method, test.code, test.response)
		})
	}
}
