package timerengine

import (
	"sync"

	"go.uber.org/zap"
)

// TaskFunc define the task func type
type TaskFunc func()

// Run adapt the TaskFunc to Task interface
func (f TaskFunc) Run() {
	f()
}

// TimerTask define the timer task
type TimerTask struct {
	sync.Mutex
	ID        uint64
	taskFunc  TaskFunc
	cycleNum  int
	engine    *TimerEngine
	started   bool
	cancelled bool
}

func newTimerTask(engine *TimerEngine, cycleNum int, f TaskFunc) *TimerTask {
	taskID := engine.nextTaskID()

	task := &TimerTask{
		ID:       taskID,
		taskFunc: f,
		cycleNum: cycleNum,
		engine:   engine,
	}
	return task
}

// Cancel cancel the task
func (task *TimerTask) Cancel() (ok bool) {
	task.Lock()
	if !task.started {
		task.cancelled = true
	}
	ok = task.cancelled
	task.Unlock()
	return
}

func (task *TimerTask) ready() (ready bool) {
	task.Lock()

	if task.cancelled {
		ready = true
	} else {
		task.cycleNum--

		if task.cycleNum <= 0 {
			ready = true
		}
	}
	task.Unlock()
	return
}

func (task *TimerTask) dispose() {
	var ok bool

	task.Lock()
	if !task.cancelled {
		task.started = true
	}
	ok = task.started
	task.Unlock()

	if ok {
		task.engine.execute(func() {
			defer func() {
				if r := recover(); r != nil {
					task.engine.logger.Error("got panic", zap.Any("error", r), zap.Stack("stack"))
				}
			}()

			if task.taskFunc != nil {
				task.taskFunc()
			}
		})
	}
}
