package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/mrapry/go-lib/codebase/interfaces"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/store"

	"github.com/gomodule/redigo/redis"
)

type redisInstance struct {
	read, write *redis.Pool
	store       interfaces.Store
}

func (m *redisInstance) ReadPool() *redis.Pool {
	return m.read
}

func (m *redisInstance) WritePool() *redis.Pool {
	return m.write
}

func (m *redisInstance) Store() interfaces.Store {
	return m.store
}

func (m *redisInstance) Disconnect(ctx context.Context) (err error) {
	deferFunc := logger.LogWithDefer("redis: disconnect...")
	defer deferFunc()

	if err := m.read.Close(); err != nil {
		return err
	}
	return m.write.Close()
}

// InitRedis connection
func InitRedis() interfaces.RedisPool {
	deferFunc := logger.LogWithDefer("Load Redis connection...")
	defer deferFunc()

	inst := new(redisInstance)

	hostRead, portRead, passRead := os.Getenv("REDIS_READ_HOST"), os.Getenv("REDIS_READ_PORT"), os.Getenv("REDIS_READ_AUTH")
	tlsRead, _ := strconv.ParseBool(os.Getenv("REDIS_READ_TLS"))
	hostWrite, portWrite, passWrite := os.Getenv("REDIS_WRITE_HOST"), os.Getenv("REDIS_WRITE_PORT"), os.Getenv("REDIS_WRITE_AUTH")
	tlsWrite, _ := strconv.ParseBool(os.Getenv("REDIS_WRITE_TLS"))

	inst.read = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", hostRead, portRead), redis.DialPassword(passRead), redis.DialUseTLS(tlsRead))
		},
	}

	pingRead := inst.read.Get()
	defer pingRead.Close()
	_, err := pingRead.Do("PING")
	if err != nil {
		panic("redis read: " + err.Error())
	}

	inst.write = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", hostWrite, portWrite), redis.DialPassword(passWrite), redis.DialUseTLS(tlsWrite))
		},
	}

	pingWrite := inst.write.Get()
	defer pingWrite.Close()
	_, err = pingWrite.Do("PING")
	if err != nil {
		panic("redis write: " + err.Error())
	}

	inst.store = store.NewRedisStore(inst.read, inst.write)

	return inst
}
