package cronworker

// cron scheduler worker, create with 100% pure internal go library (using reflect select channel)

import (
	"context"
	"fmt"
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
	service    factory.ServiceFactory
	runningJob int
	shutdown   chan struct{}
	wg         sync.WaitGroup
}

// NewWorker create new cron worker
func NewWorker(service factory.ServiceFactory) factory.AppServerFactory {
	refreshWorkerNotif = make(chan struct{})
	return &cronWorker{
		service:  service,
		shutdown: make(chan struct{}),
	}
}

func (c *cronWorker) Serve() {

	// add shutdown channel to first index
	workers = append(workers, reflect.SelectCase{
		Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c.shutdown),
	})
	// add refresh worker channel to second index
	workers = append(workers, reflect.SelectCase{
		Dir: reflect.SelectRecv, Chan: reflect.ValueOf(refreshWorkerNotif),
	})

	for _, m := range c.service.GetModules() {
		if h := m.WorkerHandler(types.Scheduler); h != nil {
			for topic, handler := range h.MountHandlers() {

				funcName, interval := golibhelper.ParseCronJobKey(topic)

				var job Job
				job.HandlerName = funcName
				job.HandlerFunc = handler
				job.Interval = interval
				if err := AddJob(job); err != nil {
					panic(err)
				}

				logger.LogYellow(fmt.Sprintf(`[CRON-WORKER] job_name: %s~%s -> every: %s`, m.Name(), funcName, interval))
			}
		}
	}

	fmt.Printf("\x1b[34;1m⇨ Cron worker running with %d jobs\x1b[0m\n\n", len(activeJobs))
	for {
		chosen, _, ok := reflect.Select(workers)
		if !ok {
			continue
		}

		// if shutdown channel captured, break loop (no more jobs will run)
		if chosen == 0 {
			break
		}

		// notify for refresh worker
		if chosen == 1 {
			continue
		}

		chosen = chosen - 2
		job := activeJobs[chosen]
		if job.nextDuration != nil {
			job.ticker.Stop()
			job.ticker = time.NewTicker(*job.nextDuration)
			workers[job.WorkerIndex].Chan = reflect.ValueOf(job.ticker.C)
			activeJobs[chosen].nextDuration = nil
		}

		c.wg.Add(1)
		c.runningJob++
		go func(job *Job) {
			defer c.wg.Done()

			trace := tracer.StartTrace(context.Background(), "CronScheduler")
			defer trace.Finish()
			ctx := trace.Context()

			defer func() {
				if r := recover(); r != nil {
					trace.SetError(fmt.Errorf("%v", r))
				}
				c.runningJob--
				logger.LogGreen(tracer.GetTraceURL(ctx))
			}()

			tags := trace.Tags()
			tags["jobName"] = job.HandlerName
			if err := job.HandlerFunc(ctx, []byte(job.Params)); err != nil {
				panic(err)
			}
		}(job)

	}
}

func (c *cronWorker) Shutdown(ctx context.Context) {
	deferFunc := logger.LogWithDefer("Stopping cron job scheduler worker...")
	defer deferFunc()

	if len(activeJobs) == 0 {
		return
	}

	c.shutdown <- struct{}{}

	done := make(chan struct{})
	go func() {
		if c.runningJob != 0 {
			fmt.Printf("\ncronjob: waiting %d job... ", c.runningJob)
		}
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
