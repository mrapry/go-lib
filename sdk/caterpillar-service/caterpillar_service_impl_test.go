package caterpillar_service

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/stretchr/testify/assert"
)

var (
	urlMock = "http://caterpillar.pln.tests"
)

func TestNewCaterpillarService(t *testing.T) {
	testName := "Test new caterpillar service"

	t.Run(testName, func(t *testing.T) {
		// set service
		service := NewCaterpillarService(urlMock)

		assert.NotNil(t, service)
	})
}

func TestCaterpillarServiceImplGetDataUnitByCode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := map[string]struct {
		wantError bool
		wantMock  bool
		status    int
		response  interface{}
	}{
		"Test #1 positive caterpillar service get data unit by code": {
			wantError: false,
			wantMock:  true,
			response:  map[string]interface{}{"success": true, "message": "success", "code": http.StatusOK},
			status:    http.StatusOK,
		},
		"Test #2 negative caterpillar service get data unit by code": {
			wantError: true,
			wantMock:  false,
			response:  map[string]interface{}{"success": false, "message": golibshared.ErrorGeneral, "code": http.StatusBadGateway},
			status:    http.StatusBadGateway,
		},
		"Test #3 negative caterpillar service get data unit by code unmarshal data": {
			wantError: true,
			wantMock:  true,
			response:  map[string]interface{}{"message": false},
			status:    http.StatusBadGateway,
		},
	}

	// set service
	service := NewCaterpillarService(urlMock)

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			unitCode := gofakeit.Digit()
			mockUrl := fmt.Sprintf("%s/v1/unit/%s", urlMock, unitCode)

			if test.wantMock {
				golibshared.CreateHttpRequestMock(http.MethodGet, mockUrl, test.status, test.response)
			} else {
				httpmock.Reset()
			}

			err := <-service.GetDataUnit(context.Background(), unitCode)
			if !test.wantError {
				assert.NoError(t, err.Error)
			} else {
				assert.Error(t, err.Error)
			}
		})
	}
}
