package golibutil

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	testName := "Test positive new http request"

	t.Run(testName, func(t *testing.T) {
		// set new request
		request := NewHTTPRequest(1, 500*time.Millisecond)

		assert.NotNil(t, request)
	})
}

func TestRequestDo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock data
	urlMock := "http://pln.tests"
	headerMock := map[string]string{"Content-Type": "application/json"}
	successResponseMock := map[string]interface{}{"success": true, "message": "success", "code": http.StatusOK}
	errorResponseMock := map[string]interface{}{"success": false, "message": "error", "code": http.StatusBadGateway}

	testCase := map[string]struct {
		wantError bool
		url       string
		code      int
		method    string
		body      interface{}
		header    map[string]string
		response  interface{}
	}{
		"Test #1 positive http request do": {
			wantError: false,
			url:       urlMock,
			code:      http.StatusOK,
			method:    http.MethodPost,
			response:  successResponseMock,
			body:      nil,
			header:    headerMock,
		},
		"Test #2 negative http request do client request": {
			wantError: true,
			url:       urlMock,
			code:      http.StatusBadGateway,
			method:    http.MethodPut,
			response:  errorResponseMock,
			body: &golibshared.Result{
				Data: false,
			},
			header: map[string]string{},
		},
	}

	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			// set new request
			request := NewHTTPRequest(1, 1*time.Millisecond)

			if test.code < 500 {
				golibshared.CreateHttpRequestMock(http.MethodPost, test.url, test.code, test.response)
			}

			var (
				err error
			)

			if test.body != nil {
				req, _ := json.Marshal(test.body)

				// do request
				_, err = request.Do(context.Background(), test.method, test.url, req, test.header)
			} else {
				// do request
				_, err = request.Do(context.Background(), test.method, test.url, nil, test.header)
			}

			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
