package cronworker

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/mrapry/go-lib/codebase/factory/types"
)

// Job model
type Job struct {
	HandlerName  string                  `json:"handler_name"`
	Interval     string                  `json:"interval"`
	HandlerFunc  types.WorkerHandlerFunc `json:"-"`
	Params       string                  `json:"params"`
	WorkerIndex  int                     `json:"worker_index"`
	ticker       *time.Ticker
	nextDuration *time.Duration
}

var activeJobs []*Job
var workers []reflect.SelectCase
var refreshWorkerNotif chan struct{}
var mutex sync.Mutex

// GetActiveJobs get registered jobs
func GetActiveJobs() []*Job {
	return activeJobs
}

// UpdateIntervalActiveJob update active job
func UpdateIntervalActiveJob(jobNumber int, newInterval string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	mutex.Lock()
	defer mutex.Unlock()

	job := activeJobs[jobNumber]
	job.Interval = newInterval
	job.nextDuration = nil

	duration, errParse := time.ParseDuration(newInterval)
	if errParse != nil {
		durationParser, nextDuration, errParse := parseAtTime(newInterval)
		if errParse != nil {
			return errParse
		}
		job.nextDuration = &nextDuration
		duration = durationParser
	}

	job.ticker.Stop()
	job.ticker = time.NewTicker(duration)
	workers[job.WorkerIndex].Chan = reflect.ValueOf(job.ticker.C)
	refreshWorkerNotif <- struct{}{}

	return
}

// AddJob to cron worker
func AddJob(job Job) error {
	mutex.Lock()
	defer mutex.Unlock()

	if job.HandlerFunc == nil {
		return errors.New("handler func cannot nil")
	}
	if job.HandlerName == "" {
		return errors.New("handler name cannot empty")
	}

	duration, err := time.ParseDuration(job.Interval)
	if err != nil {
		durationParser, nextDuration, err := parseAtTime(job.Interval)
		if err != nil {
			return err
		}
		job.nextDuration = &nextDuration
		duration = durationParser
	}

	job.ticker = time.NewTicker(duration)
	job.WorkerIndex = len(workers)

	activeJobs = append(activeJobs, &job)
	workers = append(workers, reflect.SelectCase{
		Dir: reflect.SelectRecv, Chan: reflect.ValueOf(job.ticker.C),
	})

	return nil
}
