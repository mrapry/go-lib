package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mrapry/go-lib/codebase/interfaces"
	"github.com/mrapry/go-lib/logger"
)

type postgreInstance struct {
	read, write *gorm.DB
}

func (m *postgreInstance) ReadDB() *gorm.DB {
	return m.read
}
func (m *postgreInstance) WriteDB() *gorm.DB {
	return m.write
}
func (m *postgreInstance) Disconnect(ctx context.Context) (err error) {
	deferFunc := logger.LogWithDefer("postgredb: disconnect...")
	defer deferFunc()

	if err := m.write.DB().Close(); err != nil {
		return err
	}
	return m.read.DB().Close()
}

// InitPosgreDB return mongo db read & write instance
func InitPosgreDB(ctx context.Context) interfaces.PostgreDatabase {
	deferFunc := logger.LogWithDefer("Load PostgreDB connection...")
	defer deferFunc()

	// create db instance
	dbInstance := new(postgreInstance)
	dbName, ok := os.LookupEnv("SQL_DATABASE_NAME")
	if !ok {
		panic("missing SQL_DATABASE_NAME environment")
	}
	// get write postgre from env
	hostRead := os.Getenv("SQL_DB_READ_HOST")
	hostWrite := os.Getenv("SQL_DB_WRITE_HOST")

	// connect to postgreDB
	read, err := gorm.Open(os.Getenv("SQL_DRIVER_NAME"),
		"host="+hostRead+
			" port="+os.Getenv("SQL_DB_READ_PORT")+
			" user="+os.Getenv("SQL_DB_READ_USER")+
			" dbname="+dbName+
			" password="+os.Getenv("DB_READ_PASS")+
			" sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("postgree: %v, conn: %s", err, hostRead))
	}
	dbInstance.read = read

	// connect to postgreDB
	write, err := gorm.Open(os.Getenv("SQL_DRIVER_NAME"),
		"host="+hostWrite+
			" port="+os.Getenv("SQL_DB_WRITE_PORT")+
			" user="+os.Getenv("SQL_DB_WRITE_USER")+
			" dbname="+dbName+
			" password="+os.Getenv("DB_WRITE_PASS")+
			" sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("postgree: %v, conn: %s", err, hostWrite))
	}
	dbInstance.write = write

	return dbInstance
}
