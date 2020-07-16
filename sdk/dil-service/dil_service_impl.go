package dil_service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mrapry/go-lib/golibshared"
	"github.com/mrapry/go-lib/golibutil"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/tracer"
	"github.com/mrapry/go-lib/wrapper"
	"go.uber.org/zap/zapcore"
)

type dilServiceImpl struct {
	host        string
	httpRequest golibutil.HTTPRequest
}

// NewAuthService constructor
func NewDilService(host string) ServiceDil {
	return &dilServiceImpl{
		host:        host,
		httpRequest: golibutil.NewHTTPRequest(5, 100*time.Millisecond),
	}
}

func (a *dilServiceImpl) GetDataUnit(ctx context.Context, code string) <-chan golibshared.Result {
	opName := "dil_service.get_data_unit"

	output := make(chan golibshared.Result)

	go tracer.WithTraceFunc(ctx, opName, func(c context.Context, tag map[string]interface{}) {
		defer close(output)

		var (
			response wrapper.HTTPResponse
			url      = fmt.Sprintf("%s/v1/unit/%s", a.host, code)
		)

		// http request
		resp, err := a.httpRequest.Do(c, http.MethodGet, url, nil, nil)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "http_request")
			output <- golibshared.Result{Error: err}
			return
		}

		// unmarshal to our target
		err = json.Unmarshal(resp, &response)
		if err != nil {
			logger.Log(zapcore.ErrorLevel, err.Error(), opName, "unmarshal_dil_response")
			output <- golibshared.Result{Error: err}
			return
		}

		output <- golibshared.Result{Data: response}
	})

	return output
}
