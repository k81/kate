package timerengine

import (
	"sync"

	"go.uber.org/zap"
)

// Task define the Task interface runned by timer engine
type Task interface {
	Run()
}

// TaskFunc define the task func type
type TaskFunc func()

// Run adapt the TaskFunc to taskengine.Task interface
func (f TaskFunc) Run() {
	f()
}

// TimerTask define the timer task
type TimerTask struct {
	sync.Mutex
	ID        uint64
	task      Task
	cycleNum  int
	engine    *TimerEngine
	started   bool
	cancelled bool
}

func newTimerTask(engine *TimerEngine, cycleNum int, task Task) *TimerTask {
	taskID := engine.nextTaskID()

	timerTask := &TimerTask{
		ID:       taskID,
		task:     task,
		cycleNum: cycleNum,
		engine:   engine,
	}
	return timerTask
}

// Cancel cancel the task
func (timerTask *TimerTask) Cancel() (ok bool) {
	timerTask.Lock()
	if !timerTask.started {
		timerTask.cancelled = true
	}
	ok = timerTask.cancelled
	timerTask.Unlock()
	return
}

func (timerTask *TimerTask) ready() (ready bool) {
	timerTask.Lock()

	if timerTask.cancelled {
		ready = true
	} else {
		timerTask.cycleNum--

		if timerTask.cycleNum <= 0 {
			ready = true
		}
	}
	timerTask.Unlock()
	return
}

func (timerTask *TimerTask) dispose() {
	var ok bool

	timerTask.Lock()
	if !timerTask.cancelled {
		timerTask.started = true
	}
	ok = timerTask.started
	timerTask.Unlock()

	if ok {
		timerTask.engine.execute(func() {
			defer func() {
				if r := recover(); r != nil {
					timerTask.engine.logger.Error("got panic", zap.Any("error", r), zap.Stack("stack"))
				}
			}()

			if timerTask.task != nil {
				timerTask.task.Run()
			}
		})
	}
}
