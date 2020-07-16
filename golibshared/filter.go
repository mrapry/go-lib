package golibshared

import (
	"github.com/Kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

const (
	SortDesc       = "desc"
	SortAsc        = "asc"
	DefaultOrderBy = "created_at"
)

// Filter basic filter model
type Filter struct {
	Limit   int32  `json:"limit" default:"10"`
	Page    int32  `json:"page" default:"1"`
	Offset  int32  `json:"-"`
	Search  string `json:"search,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
	Sort    string `json:"sort,omitempty" default:"desc" lower:"true"`
	SortInt int    `json:"-"`
	ShowAll bool   `json:"showAll" default:"false"`
}

// SetSort method
func (f *Filter) SetSort() {
	switch strings.ToLower(f.Sort) {
	case SortAsc:
		f.SortInt = 1
	case SortDesc:
		f.SortInt = -1
	}
}

// SetOrderBy method
func (f *Filter) SetOrderBy(fieldMap map[string]string) bson.M {
	var orderBy = bson.M{DefaultOrderBy: f.SortInt}

	if f.OrderBy != "" {
		if orderBy, ok := fieldMap[f.OrderBy]; ok {
			f.OrderBy = orderBy
		}
		orderBy = bson.M{f.OrderBy: f.SortInt}
	}

	return orderBy
}

// SetSearch method
func (f *Filter) SetSearch(query []bson.M, fields []string) []bson.M {
	var where []bson.M

	if f.Search != "" {
		for _, field := range fields {
			where = append(where, bson.M{
				field: bson.M{operator.Regex: primitive.Regex{Pattern: "^" + f.Search + ".*", Options: "i"}},
			},
			)
		}
	}

	if len(where) > 0 {
		query = append(query, bson.M{operator.Or: where})
	}

	return query
}

// CalculateOffset method
func (f *Filter) CalculateOffset() {
	f.Offset = (f.Page - 1) * f.Limit
}
