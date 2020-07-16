package tracer

import (
	"context"
	"github.com/brianvoe/gofakeit/v5"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestTraceMongoSetTags(t *testing.T) {
	testName := "Test trace mongo set tags"

	t.Run(testName, func(t *testing.T) {
		traceMongo := TraceMongo{
			Collection: gofakeit.Word(),
			Method:     Find,
			Filter:     bson.M{gofakeit.Word(): gofakeit.Word()},
			Sort:       bson.M{gofakeit.Word(): -1},
			Skip:       gofakeit.Int64(),
			Limit:      gofakeit.Int64(),
		}
		traceMongo.SetTags(context.Background())
	})
}
