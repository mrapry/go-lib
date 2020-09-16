package golibhelper

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFromQueryParam(t *testing.T) {
	type Embed struct {
		Page   int    `json:"page"`
		Offset int    `json:"-"`
		Sort   string `json:"sort,omitempty" default:"desc" lower:"true"`
	}
	type params struct {
		Embed
		IsActive bool    `json:"isActive"`
		Ptr      *string `json:"ptr"`
	}

	t.Run("Testcase #1: Positive", func(t *testing.T) {
		urlVal, err := url.ParseQuery("page=1&ptr=val&isActive=true")
		assert.NoError(t, err)

		var p params
		err = ParseFromQueryParam(urlVal, &p)
		assert.NoError(t, err)
		assert.Equal(t, p.Page, 1)
		assert.Equal(t, *p.Ptr, "val")
		assert.Equal(t, p.IsActive, true)
	})
	t.Run("Testcase #2: Negative, invalid data type (string to int in struct)", func(t *testing.T) {
		urlVal, err := url.ParseQuery("page=undefined")
		assert.NoError(t, err)

		var p params
		err = ParseFromQueryParam(urlVal, &p)
		assert.Error(t, err)
	})
	t.Run("Testcase #3: Negative, invalid data type (not boolean)", func(t *testing.T) {
		urlVal, err := url.ParseQuery("isActive=terue")
		assert.NoError(t, err)

		var p params
		err = ParseFromQueryParam(urlVal, &p)
		assert.Error(t, err)
	})
	t.Run("Testcase #4: Negative, invalid target type (not pointer)", func(t *testing.T) {
		urlVal, err := url.ParseQuery("isActive=true")
		assert.NoError(t, err)

		var p params
		err = ParseFromQueryParam(urlVal, p)
		assert.Error(t, err)
	})
}
