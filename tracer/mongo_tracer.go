package tracer

import (
	"context"
	"fmt"
)

const (
	// db tag error
	DBType       = "db.type"
	DBMethod     = "db.method"
	DBCollection = "db.collection"
	DBStatement  = "db.statement"
	DBSort       = "db.sort"
	DBSkip       = "db.skip"
	DBLimit      = "db.limit"

	// db method for tracer
	Aggregate     = "Aggregate"
	BulkWrite     = "BulkWrite"
	CountDocument = "CountDocument"
	Find          = "Find"
	FindOne       = "FindOne"
	DeleteMany    = "DeleteMany"
	DeleteOne     = "DeleteOne"
	Insert        = "Insert"
	InsertMany    = "InsertMany"
	InsertOne     = "InsertOne"
	UpdateOne     = "UpdateOne"
	UpdateMany    = "UpdateMany"
)

type TraceMongo struct {
	Collection string
	Method     string
	Filter     interface{}
	Sort       interface{}
	Skip       int64
	Limit      int64
}

// SetTags function from tags, filter must be bson.M for a  better result
func (t *TraceMongo) SetTags(ctx context.Context) {
	// init tracer
	tracer := StartTrace(ctx, fmt.Sprintf("mongodb:%s", t.Method))
	defer tracer.Finish()

	// set tags
	tags := tracer.Tags()

	// set data
	tags[DBType] = "mongo"
	tags[DBMethod] = t.Method

	if t.Collection != "" {
		tags[DBCollection] = t.Collection
	}

	if t.Filter != nil {
		tags[DBStatement] = t.Filter
	}

	if t.Sort != nil {
		tags[DBSort] = t.Sort
	}

	if t.Limit > 0 {
		tags[DBLimit] = t.Limit
	}

	if t.Skip > 0 {
		tags[DBSkip] = t.Skip
	}
}
