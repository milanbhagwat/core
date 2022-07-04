package runner

import (
	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/support/log"
	"sync"
	"time"
)

func NewRunnerTracker() *RunnerTracker {
	return &RunnerTracker{runnertrackerwg: &sync.WaitGroup{}}
}

type RunnerTracker struct {
	runnertrackerwg *sync.WaitGroup
}

func (rt RunnerTracker) AddTracker() {
	rt.runnertrackerwg.Add(1)
}

func (rt RunnerTracker) DoneTracker() {
	rt.runnertrackerwg.Done()
}

func (rt RunnerTracker) WaitForTrackerAllDone() {
	rt.runnertrackerwg.Wait()
}

func (rt RunnerTracker) WaitForActionsCompletion(timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		rt.WaitForTrackerAllDone()
	}()
	select {
	case <-c:
		return false // actions completed
	case <-time.After(timeout):
		return true // timed out
	}
}

func (rt RunnerTracker) gracefulStop() {
	logger := log.RootLogger()
	delayedStopInterval := app.GetDelayedStopInterval()
	if delayedStopInterval != "" {
		duration, err := time.ParseDuration(delayedStopInterval)
		if err != nil {
			logger.Errorf("Invalid interval - %s  specified for delayed stop. It must suffix with time unit e.g. %sms, %ss", delayedStopInterval, delayedStopInterval, delayedStopInterval)
		} else {
			logger.Infof("Delaying application stop by max - %s", delayedStopInterval)
			if isTimeout := rt.WaitForActionsCompletion(duration); isTimeout {
				logger.Info("All actions not completed before engine shutdown")
			} else {
				logger.Info("All actions completed before engine shutdown")
			}
		}

	}
}
