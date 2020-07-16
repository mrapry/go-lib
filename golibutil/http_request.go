package golibutil

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
	"github.com/mrapry/go-lib/tracer"
)

// Request struct
type Request struct {
	client *httpclient.Client
}

// HTTPRequest interface
type HTTPRequest interface {
	Do(context context.Context, method, url string, reqBody []byte, headers map[string]string) ([]byte, error)
}

// NewHTTPRequest function
// Request's Constructor
// Returns : *Request
func NewHTTPRequest(retries int, sleepBetweenRetry time.Duration) HTTPRequest {
	// define a maximum jitter interval
	maximumJitterInterval := 5 * time.Millisecond

	// create a backoff
	backoff := heimdall.NewConstantBackoff(sleepBetweenRetry, maximumJitterInterval)

	// create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	// set http timeout
	timeout := 10000 * time.Millisecond

	// set http client
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(retries),
	)

	return &Request{
		client: client,
	}
}

// Do function, for http client call
func (request *Request) Do(context context.Context, method, url string, requestBody []byte, headers map[string]string) ([]byte, error) {
	var (
		respBody   []byte
		respStatus string
	)

	// set request http
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// iterate optional data of headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// set tracer
	trace := tracer.StartTrace(context, fmt.Sprintf("%s %s%s", method, req.URL.Host, req.URL.Path))
	trace.InjectHTTPHeader(req)
	defer func() {
		tags := map[string]interface{}{
			"http.headers":    req.Header,
			"http.method":     req.Method,
			"http.url":        req.URL.String(),
			"response.status": respStatus,
			"response.body":   string(respBody),
		}

		if requestBody != nil {
			tags["request.body"] = string(requestBody)
		}

		trace.Finish(tags)
	}()

	// client request
	r, err := request.client.Do(req)
	if err != nil {
		return nil, err
	}
	// close response body
	defer r.Body.Close()

	respBody, err = ioutil.ReadAll(r.Body)
	return respBody, err
}
