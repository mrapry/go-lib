package redisworker

// Redis subscriber worker codebase

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/mrapry/go-lib/codebase/factory"
	"github.com/mrapry/go-lib/codebase/factory/types"
	"github.com/mrapry/go-lib/golibhelper"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/tracer"
)

type redisWorker struct {
	pubSubConn func() redis.PubSubConn
	isHaveJob  bool
	service    factory.ServiceFactory
	shutdown   chan struct{}
	wg         sync.WaitGroup
}

// NewWorker create new redis subscriber
func NewWorker(service factory.ServiceFactory) factory.AppServerFactory {
	redisPool := service.GetDependency().GetRedisPool().WritePool()

	return &redisWorker{
		service: service,
		pubSubConn: func() redis.PubSubConn {
			conn := redisPool.Get()
			conn.Do("CONFIG", "SET", "notify-keyspace-events", "Ex")

			psc := redis.PubSubConn{Conn: conn}
			psc.PSubscribe("__keyevent@*__:expired")

			return psc
		},
		shutdown: make(chan struct{}),
	}
}

func (r *redisWorker) Serve() {
	handlers := make(map[string]types.WorkerHandlerFunc)
	for _, m := range r.service.GetModules() {
		if h := m.WorkerHandler(types.RedisSubscriber); h != nil {
			for topic, handlerFunc := range h.MountHandlers() {
				fmt.Println(golibhelper.StringYellow(fmt.Sprintf(`[REDIS-SUBSCRIBER] (key prefix): "%-10s"  (processed by module): %s`, topic, m.Name())))
				handlers[golibhelper.BuildRedisPubSubKeyTopic(string(m.Name()), topic)] = handlerFunc
			}
		}
	}

	if len(handlers) == 0 {
		log.Println("redis subscriber: no topic provided")
		return
	}
	r.isHaveJob = true

	psc := r.pubSubConn()

	// listen redis subscriber
	messageReceiver := make(chan []byte)
	go func() {
		for {
			switch msg := psc.Receive().(type) {
			case redis.Message:
				messageReceiver <- msg.Data
			case error:
				// if network connection error, create new connection from pool
				psc = r.pubSubConn()
			}
		}
	}()

	// run worker with listen shutdown channel
	fmt.Printf("\x1b[34;1m⇨ Redis pubsub worker running with %d keys\x1b[0m\n\n", len(handlers))
	for {
		select {
		case message := <-messageReceiver:
			r.wg.Add(1)
			go tracer.WithTraceFunc(context.Background(), "RedisSubscriber", func(ctx context.Context, tags map[string]interface{}) {
				defer r.wg.Done()
				defer func() {
					if r := recover(); r != nil {
						tracer.SetError(ctx, fmt.Errorf("%v", r))
					}
				}()
				tags["message"] = string(message)

				handlerName, messageData := golibhelper.ParseRedisPubSubKeyTopic(string(message))
				handlerFunc, ok := handlers[handlerName]
				if !ok {
					return
				}
				if err := handlerFunc(ctx, []byte(messageData)); err != nil {
					panic(err)
				}
			})
		case <-r.shutdown:
			break
		}
	}
}

func (r *redisWorker) Shutdown(ctx context.Context) {
	deferFunc := logger.LogWithDefer("Stopping redis subscriber worker...")
	defer deferFunc()

	if !r.isHaveJob {
		return
	}

	r.shutdown <- struct{}{}

	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		fmt.Print("redis-subscriber: force shutdown ")
	case <-done:
		fmt.Print("redis-subscriber: success waiting all job until done ")
	}
}
