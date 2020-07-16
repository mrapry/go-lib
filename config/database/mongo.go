package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Kamva/mgm/v3"
	"github.com/mrapry/go-lib/codebase/interfaces"
	"github.com/mrapry/go-lib/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoInstance struct {
	read, write *mongo.Database
}

func (m *mongoInstance) ReadDB() *mongo.Database {
	return m.read
}
func (m *mongoInstance) WriteDB() *mongo.Database {
	return m.write
}
func (m *mongoInstance) Disconnect(ctx context.Context) (err error) {
	deferFunc := logger.LogWithDefer("mongodb: disconnect...")
	defer deferFunc()

	if err := m.write.Client().Disconnect(ctx); err != nil {
		return err
	}
	return m.read.Client().Disconnect(ctx)
}

// InitMongoDB return mongo db read & write instance
func InitMongoDB(ctx context.Context) interfaces.MongoDatabase {
	deferFunc := logger.LogWithDefer("Load MongoDB connection...")
	defer deferFunc()

	// create db instance
	dbInstance := new(mongoInstance)
	dbName, ok := os.LookupEnv("MONGODB_DATABASE_NAME")
	if !ok {
		panic("missing MONGODB_DATABASE_NAME environment")
	}

	// set default mgm write
	mgm.SetDefaultConfig(&mgm.Config{CtxTimeout: 15000 * time.Millisecond}, dbName)

	// get write mongo from env
	hostWrite := os.Getenv("MONGODB_HOST_WRITE")

	// connect to MongoDB
	client, err := mgm.NewClient(options.Client().ApplyURI(hostWrite))
	if err != nil {
		panic(fmt.Errorf("mongo: %v, conn: %s", err, hostWrite))
	}
	dbInstance.write = client.Database(dbName)

	// get read mongo from env
	hostRead := os.Getenv("MONGODB_HOST_READ")

	// connect to MongoDB
	client, err = mgm.NewClient(options.Client().ApplyURI(hostRead))
	if err != nil {
		panic(fmt.Errorf("mongo: %v, conn: %s", err, hostRead))
	}
	dbInstance.read = client.Database(dbName)

	return dbInstance
}
