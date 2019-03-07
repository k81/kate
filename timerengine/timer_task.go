package timerengine

import (
	"context"
	"sync"

	"github.com/k81/kate/utils"
	"github.com/k81/log"
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
	logger    *log.Logger
	ctx       context.Context
}

func newTimerTask(engine *TimerEngine, cycleNum int, f TaskFunc) *TimerTask {
	taskId := engine.nextTaskId()

	task := &TimerTask{
		Id:       taskId,
		taskFunc: f,
		cycleNum: cycleNum,
		engine:   engine,
		logger:   engine.logger.With("task_id", taskId),
		ctx:      engine.ctx,
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
			if task.logger.Enabled(log.LevelTrace) {
				task.logger.Trace(task.ctx, "timer task started")
			}

			defer func() {
				if r := recover(); r != nil {
					task.logger.Error(task.ctx, "got panic", "error", r, "stack", utils.GetPanicStack())
				}

				if task.logger.Enabled(log.LevelTrace) {
					task.logger.Trace(task.ctx, "timer task stopped")
				}
			}()

			if task.taskFunc != nil {
				task.taskFunc()
			}
		})
	}
}
