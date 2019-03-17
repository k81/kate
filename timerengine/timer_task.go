package timerengine

import (
	"context"
	"sync"

	"github.com/k81/kate/utils"
	"github.com/k81/log"
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
	ctx       context.Context
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
		ctx:      log.WithContext(engine.ctx, "task_id", taskID),
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
			if log.Enabled(log.LevelTrace) {
				log.Trace(task.ctx, "timer task started")
			}

			defer func() {
				if r := recover(); r != nil {
					log.Error(task.ctx, "got panic", "error", r, "stack", utils.GetPanicStack())
				}

				if log.Enabled(log.LevelTrace) {
					log.Trace(task.ctx, "timer task stopped")
				}
			}()

			if task.taskFunc != nil {
				task.taskFunc()
			}
		})
	}
}
