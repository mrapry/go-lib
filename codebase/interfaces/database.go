package interfaces

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

// SQLDatabase abstraction
type SQLDatabase interface {
	ReadDB() *sql.DB
	WriteDB() *sql.DB
	Closer
}

// MongoDatabase abstraction
type MongoDatabase interface {
	ReadDB() *mongo.Database
	WriteDB() *mongo.Database
	Closer
}

// RedisPool abstraction
type RedisPool interface {
	ReadPool() *redis.Pool
	WritePool() *redis.Pool
	Closer
}

// PostgreDatabase abstraction
type PostgreDatabase interface {
	ReadDB() *gorm.DB
	WriteDB() *gorm.DB
	Closer
}
