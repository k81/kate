package timerengine

import (
	"context"
	"sync"

	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"
)

type TaskFunc func()

func (f TaskFunc) Run() {
	f()
}

type TimerTask struct {
	sync.Mutex
	Id        uint64
	taskFunc  TaskFunc
	cycleNum  int
	started   bool
	engine    *TimerEngine
	cancelled bool
	ctx       context.Context
}

func newTimerTask(engine *TimerEngine, cycleNum int, f TaskFunc) *TimerTask {
	taskId := engine.nextTaskId()

	task := &TimerTask{
		Id:       taskId,
		taskFunc: f,
		cycleNum: cycleNum,
		engine:   engine,
		ctx:      log.SetContext(engine.ctx, "task_id", taskId),
	}
	return task
}

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
		task.cycleNum -= 1

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
			if log.Enabled(log.TraceLevel) {
				log.Trace(task.ctx, "timer task started")
			}

			defer func() {
				if r := recover(); r != nil {
					log.Error(task.ctx, "got panic", "error", r, "stack", utils.GetPanicStack())
				}

				if log.Enabled(log.TraceLevel) {
					log.Trace(task.ctx, "timer task stopped")
				}
			}()

			if task.taskFunc != nil {
				task.taskFunc()
			}
		})
	}
}
