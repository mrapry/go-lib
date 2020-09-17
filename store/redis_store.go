package store

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mrapry/go-lib/codebase/interfaces"
	"github.com/mrapry/go-lib/tracer"
)

// RedisStore redis
type RedisStore struct {
	read, write *redis.Pool
}

// NewRedisStore constructor
func NewRedisStore(read, write *redis.Pool) interfaces.Store {
	return &RedisStore{read: read, write: write}
}

// Get method
func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
	opName := "redis:get"

	// set tracer
	tracer := tracer.StartTrace(ctx, opName)
	tag := tracer.Tags()
	tag["db.statement"] = "GET"
	tag["db.key"] = key
	defer tracer.Finish(tag)

	// set client
	cl := r.read.Get()
	defer cl.Close()

	var data string
	var err error
	data, err = redis.String(cl.Do("GET", key))
	if err != nil {
		tracer.SetError(err)
		return data, err
	}

	return data, nil
}

// GetKeys method
func (r *RedisStore) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	opName := "redis:get_keys"

	// set tracer
	tracer := tracer.StartTrace(ctx, opName)
	tag := tracer.Tags()
	tag["db.statement"] = "KEYS"
	tag["db.key"] = pattern
	defer tracer.Finish(tag)

	// set client
	cl := r.read.Get()
	defer cl.Close()

	var datas []string
	var err error
	datas, err = redis.Strings(cl.Do("KEYS", fmt.Sprintf("%s*", pattern)))
	if err != nil {
		tracer.SetError(err)
		return datas, err
	}

	return datas, nil
}

// Set method
func (r *RedisStore) Set(ctx context.Context, key string, value interface{}, expire time.Duration) (err error) {
	opName := "redis:set"

	// set tracer
	trace := tracer.StartTrace(ctx, opName)
	defer func() {
		if err != nil {
			trace.SetError(err)
		}
		trace.Finish()
	}()

	tag := trace.Tags()
	tag["db.statement"] = "SET"
	tag["db.key"] = key
	tag["db.expired"] = expire.String()

	// set client
	cl := r.write.Get()
	defer cl.Close()

	_, err = cl.Do("SET", key, value)
	if err != nil {
		return
	}

	_, err = cl.Do("EXPIRE", key, int(expire.Seconds()))
	if err != nil {
		return
	}

	return nil
}

// Exists method
func (r RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	opName := "redis:exists"

	// set tracer
	tracer := tracer.StartTrace(ctx, opName)
	tag := tracer.Tags()
	tag["db.statement"] = "EXISTS"
	tag["db.key"] = key
	defer tracer.Finish(tag)

	// set client
	cl := r.read.Get()
	defer cl.Close()

	var ok bool
	var err error
	ok, err = redis.Bool(cl.Do("EXISTS", key))
	if err != nil {
		tracer.SetError(err)
		return ok, err
	}

	return ok, nil
}

// Delete method
func (r *RedisStore) Delete(ctx context.Context, key string) error {
	opName := "redis:delete"

	// set tracer
	tracer := tracer.StartTrace(ctx, opName)
	tag := tracer.Tags()
	tag["db.statement"] = "DEL"
	tag["db.key"] = key
	defer tracer.Finish(tag)

	// set client
	cl := r.write.Get()
	defer cl.Close()

	_, err := cl.Do("DEL", key)
	if err != nil {
		tracer.SetError(err)
		return err
	}

	return nil
}
