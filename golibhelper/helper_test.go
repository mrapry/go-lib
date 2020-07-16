package golibhelper

import (
	"github.com/brianvoe/gofakeit/v5"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestParseFromQueryParam(t *testing.T) {
	type Filter struct {
		TestString  string  `json:"test_string" lower:"true"`
		TestInt     int32   `json:"test_int"`
		TestBool    bool    `json:"test_bool"`
		TestPointer *string `json:"test_pointer"`
		TestEmpty   *string `json:"-,omitempty"`
		TestDefault *string `json:"test_default" default:"test"`
	}

	testCases := map[string]struct {
		wantError   bool
		testString  string
		testInt     string
		testBool    string
		testPointer string
		testDefault string
	}{
		"Test positive parse from query param": {
			wantError:   false,
			testString:  gofakeit.Word(),
			testInt:     cast.ToString(1),
			testBool:    cast.ToString(true),
			testPointer: gofakeit.Word(),
		},
		"negative parse from query param": {
			wantError:   true,
			testString:  gofakeit.Word(),
			testInt:     gofakeit.Word(),
			testBool:    gofakeit.Word(),
			testPointer: gofakeit.Word(),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			// create url value
			url := url.Values{}

			url.Set("test_string", test.testString)
			url.Set("test_int", test.testInt)
			url.Set("test_bool", test.testBool)
			url.Set("test_pointer", test.testPointer)

			// set filter
			var filter Filter

			err := ParseFromQueryParam(url, &filter)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
