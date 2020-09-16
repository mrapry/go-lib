package golibshared

import (
	"testing"

	"github.com/Kamva/mgm/v3/operator"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestToBSON(t *testing.T) {
	nameMock := gofakeit.Name()

	testCase := map[string]struct {
		structure interface{}
		expected  bson.M
	}{
		"Test #1 positive to bson": {
			structure: struct {
				Name string `bson:"name"`
				Age  int    `bson:"age"`
			}{
				Name: nameMock,
				Age:  20,
			},
			expected: bson.M{operator.And: []bson.M{{"name": nameMock}, {"age": 20}}},
		},
		"Test #2 negative invalid tag to bson": {
			structure: struct {
				Name string `bson:"-"`
				Age  int    `bson:"-"`
			}{
				Name: nameMock,
			},
			expected: bson.M{operator.And: []bson.M{}},
		},
		"Test #3 positive with zero value": {
			structure: struct {
				Name string `bson:"name"`
				Age  int    `bson:"age"`
			}{
				Name: "",
				Age:  0,
			},
			expected: bson.M{operator.And: []bson.M{}},
		},
	}

	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			bsonEncoding := ToBSON(test.structure)

			assert.Equal(t, test.expected, bsonEncoding)
		})
	}
}
