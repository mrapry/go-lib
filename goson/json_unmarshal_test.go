package goson

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	type Slice struct {
		FieldA uint16  `json:"fieldA"`
		FieldB string  `json:"fieldB"`
		Exist  string  `json:"exist"`
		Test   float32 `json:"test"`
	}
	type Model struct {
		ID        int      `json:"id"`
		Name      string   `json:"name"`
		MustFloat *float64 `json:"mustFloat"`
		MustInt   int      `json:"mustInt"`
		Uint      uint     `json:"uint"`
		IsExist   *bool    `json:"isExist"`
		Obj       *struct {
			N       int `json:"n"`
			Testing struct {
				Ss int `json:"ss"`
			} `json:"testing"`
		} `json:"obj"`
		Slice     []Slice   `json:"slice"`
		Strings   []*string `json:"str"`
		Ints      []int     `json:"ints"`
		Bools     []bool    `json:"bools"`
		NoTag     string
		skip      string      // cannot set value to this field because unexported
		Interface interface{} `json:"interface"`
	}

	t.Run("Testcase #1: Testing Unmarshal with root is JSON Object", func(t *testing.T) {
		data := []byte(`{
				"id": "01",
				"name": "agungdp",
				"mustFloat": "2.23423",
				"mustInt": 2.23423,
				"uint": 11,
				"isExist": "true",
				"obj": {
				  "n": 2,
				  "testing": {
					"ss": "23840923849284"
				  }
				},
				"slice": [
				  {
					"fieldA": "100",
					"fieldB": 3000,
					"exist": true,
					"test": "3.14"
				  },
				  {
					"fieldA": 50000,
					"fieldB": "123000",
					"exist": 3000,
					"test": 1323.123
				  }
				],
				"str": ["a", true],
				"ints": ["2", 3],
				"bools": [1, "true", "0", true],
				"NoTag": 19283091832,
				"interface": 1
			  }`)
		var target Model
		err := Unmarshal(data, &target)
		assert.NoError(t, err)

		assert.NotNil(t, target.MustFloat)
		assert.Equal(t, uint(11), target.Uint)
		assert.Equal(t, 23840923849284, target.Obj.Testing.Ss)
		assert.Equal(t, "true", target.Slice[0].Exist)
		assert.Equal(t, "true", *target.Strings[1])
		assert.Equal(t, 2, target.Ints[0])
		assert.Equal(t, false, target.Bools[2])
		assert.Equal(t, "19283091832", target.NoTag)
		assert.NotNil(t, target.Interface)

		fmt.Printf("%+v\n\n", target)
	})

	t.Run("Testcase #2: Testing Unmarshal with root is JSON Array", func(t *testing.T) {
		data := []byte(`[
			{
				 "fieldA": 100,
				 "fieldB": "3000",
				 "exist": "true",
				 "test": 3.14
			},
			{
				 "fieldA": 50000,
				 "fieldB": "123000",
				 "exist": "3000",
				 "test": 1323.123
			}
	   ]`)
		var target []Slice
		err := Unmarshal(data, &target)
		assert.NoError(t, err)

		assert.NotNil(t, target)
		assert.Len(t, target, 2)

		fmt.Printf("%+v\n\n", target)
	})

	t.Run("Testcase #3: Testing error unmarshal (target is not pointer)", func(t *testing.T) {
		data := []byte(`{}`)
		var target Model
		err := Unmarshal(data, target)
		assert.Error(t, err)
	})

	t.Run("Testcase #4: Testing error unmarshal (invalid json format)", func(t *testing.T) {
		data := []byte(`{1: satu}`)
		var target Model
		err := Unmarshal(data, &target)
		assert.Error(t, err)
	})
}
