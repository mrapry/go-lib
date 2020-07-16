package database

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/spf13/cast"
	"os"
	"testing"
)

func TestInitMongoDB(t *testing.T) {
	mongoUrlMock := fmt.Sprintf("mongodb://%s:%s@%s:%s", gofakeit.Word(), gofakeit.Word(), gofakeit.Word(), cast.ToString(1))

	testCase := map[string]struct {
		descriptorWrite string
		descriptorRead  string
		dbname          string
	}{
		"Test #1 positive init db mongo connection": {
			descriptorRead:  mongoUrlMock,
			descriptorWrite: mongoUrlMock,
			dbname:          gofakeit.Word(),
		},
		"Test #2 negative init db mongo connection connect db read": {
			descriptorRead:  gofakeit.Word(),
			descriptorWrite: mongoUrlMock,
			dbname:          gofakeit.Word(),
		},
		"Test #3 negative init db mongo connection connect db write": {
			descriptorWrite: gofakeit.Word(),
			dbname:          gofakeit.Word(),
		},
		"Test #4 negative init db mongo connection no db name": {
		},
	}

	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			defer func() { recover() }()

			if test.dbname != "" {
				os.Setenv("MONGODB_DATABASE_NAME", test.dbname)
			} else {
				os.Unsetenv("MONGODB_DATABASE_NAME")
			}

			if test.descriptorWrite != "" {
				os.Setenv("MONGODB_HOST_WRITE", test.descriptorWrite)
			} else {
				os.Unsetenv("MONGODB_HOST_WRITE")
			}

			if test.descriptorRead != "" {
				os.Setenv("MONGODB_HOST_READ", test.descriptorRead)
			} else {
				os.Unsetenv("MONGODB_HOST_READ")
			}

			db := InitMongoDB(context.Background())
			db.ReadDB()
			db.WriteDB()
			db.Disconnect(context.Background())
		})
	}
}
