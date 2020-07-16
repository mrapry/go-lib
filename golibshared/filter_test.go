package golibshared

import (
	"github.com/Kamva/mgm/v3/operator"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFilterCalculateOffset(t *testing.T) {
	testName := "Test calculate offset"

	t.Run(testName, func(t *testing.T) {
		filter := Filter{Limit: 1, Page: 1}

		filter.CalculateOffset()

		assert.Equal(t, int32(0), filter.Offset)
	})
}

func TestFilterSetSort(t *testing.T) {
	testCase := map[string]struct {
		sort     string
		expected int
	}{
		"Test set sort asc": {
			sort:     SortAsc,
			expected: 1,
		},
		"Test set sort desc": {
			sort:     SortDesc,
			expected: -1,
		},
	}

	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			filter := Filter{Sort: test.sort}

			filter.SetSort()

			assert.Equal(t, test.expected, filter.SortInt)
		})
	}
}

func TestFilterSetOrderBy(t *testing.T) {
	testName := "Test set order by"

	t.Run(testName, func(t *testing.T) {
		filter := Filter{OrderBy: DefaultOrderBy, Sort: SortAsc}
		fieldMap := map[string]string{
			DefaultOrderBy: DefaultOrderBy,
		}

		filter.SetSort()
		orderBy := filter.SetOrderBy(fieldMap)

		assert.Equal(t, bson.M{DefaultOrderBy: 1}, orderBy)
	})
}

func TestFilterSetSearch(t *testing.T) {
	testName := "Test set search"

	t.Run(testName, func(t *testing.T) {
		search := gofakeit.Word()
		filter := Filter{Search: search}

		query := []bson.M{}
		fields := []string{"name"}

		// search
		query = filter.SetSearch(query, fields)

		// expected
		expected := []bson.M{
			bson.M{operator.Or: []bson.M{
				bson.M{"name": bson.M{operator.Regex: primitive.Regex{Pattern: `^` + search + `.*`, Options: "i"}}},
			}},
		}

		assert.Equal(t, expected, query)
	})
}
