package golibshared

import (
	"github.com/Kamva/mgm/v3/operator"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
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
			}{
				Name: nameMock,
			},
			expected: bson.M{operator.And: []bson.M{bson.M{"name": nameMock}}},
		},
		"Test #2 negative invalid tag to bson": {
			structure: struct {
				Name string `bson:"-"`
			}{
				Name: nameMock,
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
