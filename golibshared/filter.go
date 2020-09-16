package golibshared

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	OrderBy string `json:"order_by,omitempty"`
	Sort    string `json:"sort,omitempty" default:"desc" lower:"true"`
	SortInt int    `json:"-"`
	ShowAll bool   `json:"show_all" default:"false"`
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

//SearchCondition implement function for search data  with array data coloum search
func (f *Filter) SearchCondition(searchFields []string) string {
	var likes []string
	for _, key := range searchFields {
		likes = append(likes, fmt.Sprintf(`CAST(lower(%s) as text) LIKE '%%%s%%'`, key, strings.ToLower(f.Search)))
	}

	return strings.Join(likes, " OR ")
}

// LikeCondition Implements function to query equal -> WHERE column LIKE '%value%'
// searchFields types is []string will be located as column name
// value is parameter for value of LIKE the data type is String
func (f *Filter) LikeCondition(searchFields []string, value string) string {
	var likes []string
	for _, key := range searchFields {
		likes = append(likes, fmt.Sprintf(`CAST(lower(%s) as text) LIKE '%%%s%%'`, key, strings.ToLower(strings.TrimSpace(value))))
	}

	return strings.Join(likes, " OR ")
}

//DateCondition is function for searching where data type is Timestamp
//string return format is column BETWEEN '2006-01-02 15:04:05' AND '2006-01-02 15:04:05'
func (f *Filter) DateCondition(column, date string) string {
	layoutFormat := "2006-01-02"
	initial, _ := time.Parse(layoutFormat, date)

	start := initial.Format("2006-01-02 15:04:05")
	end := initial.Add(time.Hour*23 + time.Minute*59 + time.Second*59).Format("2006-01-02 15:04:05")
	getDate := fmt.Sprintf("%s BETWEEN '%s' AND '%s'", column, start, end)

	return getDate
}
