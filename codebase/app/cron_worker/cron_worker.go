package cronworker

// cron scheduler worker, create with 100% pure internal go library (using reflect select channel)

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/mrapry/go-lib/codebase/factory"
	"github.com/mrapry/go-lib/codebase/factory/types"
	"github.com/mrapry/go-lib/golibhelper"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/tracer"
)

type cronWorker struct {
	service   factory.ServiceFactory
	isHaveJob bool
	shutdown  chan struct{}
	wg        sync.WaitGroup
}

// NewWorker create new cron worker
func NewWorker(service factory.ServiceFactory) factory.AppServerFactory {
	return &cronWorker{
		service:  service,
		shutdown: make(chan struct{}),
	}
}

func (c *cronWorker) Serve() {
	var jobs []schedulerJob
	var schedulerChannels []reflect.SelectCase
	for _, m := range c.service.GetModules() {
		if h := m.WorkerHandler(types.Scheduler); h != nil {
			for topic, handler := range h.MountHandlers() {
				var job schedulerJob

				funcName, interval := golibhelper.ParseCronJobKey(topic)
				duration, err := time.ParseDuration(interval)
				if err != nil {
					durationParser, nextDuration, err := parseAtTime(interval)
					if err != nil {
						panic(err)
					}
					job.nextDuration = &nextDuration
					duration = durationParser
				}

				job.handlerName = funcName
				job.ticker = time.NewTicker(duration)
				job.handlerFunc = handler

				schedulerChannels = append(schedulerChannels, reflect.SelectCase{
					Dir: reflect.SelectRecv, Chan: reflect.ValueOf(job.ticker.C),
				})
				jobs = append(jobs, job)

				fmt.Println(golibhelper.StringYellow(fmt.Sprintf(`[CRON-WORKER] job_name: "%s" -> every: %s`, funcName, interval)))
			}
		}
	}

	if len(jobs) == 0 {
		log.Println("cronjob: no scheduler handler found")
		return
	}

	c.isHaveJob = true

	// add shutdown channel to last index
	schedulerChannels = append(schedulerChannels, reflect.SelectCase{
		Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c.shutdown),
	})

	fmt.Printf("\x1b[34;1m⇨ Cron worker running with %d jobs\x1b[0m\n\n", len(jobs))
	for {
		chosen, _, ok := reflect.Select(schedulerChannels)
		if !ok {
			continue
		}

		// if shutdown channel captured, break loop (no more jobs will run)
		if chosen == len(schedulerChannels)-1 {
			break
		}

		job := jobs[chosen]
		if job.nextDuration != nil {
			job.ticker.Stop()
			job.ticker = time.NewTicker(*job.nextDuration)
			schedulerChannels[chosen].Chan = reflect.ValueOf(job.ticker.C)
			jobs[chosen].nextDuration = nil
		}

		c.wg.Add(1)
		go func(job schedulerJob) {
			defer c.wg.Done()

			trace := tracer.StartTrace(context.Background(), "CronScheduler")
			defer trace.Finish()
			ctx := trace.Context()

			defer func() {
				if r := recover(); r != nil {
					trace.SetError(fmt.Errorf("%v", r))
				}
			}()

			tags := trace.Tags()
			tags["jobName"] = job.handlerName
			if err := job.handlerFunc(ctx, []byte(job.params)); err != nil {
				panic(err)
			}
		}(job)

	}
}

func (c *cronWorker) Shutdown(ctx context.Context) {
	deferFunc := logger.LogWithDefer("Stopping cron job scheduler worker...")
	defer deferFunc()

	if !c.isHaveJob {
		return
	}

	c.shutdown <- struct{}{}

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		fmt.Print("cronjob: force shutdown ")
	case <-done:
		fmt.Print("cronjob: success waiting all job until done ")
	}
}

type schedulerJob struct {
	ticker       *time.Ticker
	nextDuration *time.Duration
	handlerName  string
	handlerFunc  types.WorkerHandlerFunc
	params       string
}
