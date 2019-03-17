package taskengine

import (
	"fmt"
	"sync"

	"context"

	"github.com/k81/kate/utils"
	"github.com/k81/log"
)

// TaskEngine define the task engine
type TaskEngine struct {
	name              string
	concurrencyTokens chan struct{}
	ctx               context.Context
	cancel            context.CancelFunc
	shutdown          bool
	sync.WaitGroup
}

// New create a new task engine
func New(ctx context.Context, name string, concurrencyLevel int) *TaskEngine {
	newctx, cancel := context.WithCancel(ctx)
	newctx = log.WithContext(newctx, "taskengine", name)

	engine := &TaskEngine{
		name:   name,
		ctx:    newctx,
		cancel: cancel,
	}

	if concurrencyLevel > 0 {
		engine.concurrencyTokens = make(chan struct{}, concurrencyLevel)
	}
	return engine
}

// Schedule schedule a task running on engine
func (engine *TaskEngine) Schedule(task Task) bool {
	if engine.shutdown {
		log.Error(engine.ctx, "already stopped, should not schedule new task")
		return false
	}

	engine.Add(1)

	if engine.concurrencyTokens != nil {
		engine.concurrencyTokens <- struct{}{}
	}

	go engine.run(task)
	return true
}

func (engine *TaskEngine) run(task Task) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(engine.ctx, "task panic:", "error", r, "stack", utils.GetPanicStack())
		}
		if engine.concurrencyTokens != nil {
			<-engine.concurrencyTokens
		}
		engine.Done()
	}()

	task.Run()
}

// Shutdown stop the task engine
func (engine *TaskEngine) Shutdown() {
	if engine.shutdown {
		panic(fmt.Sprintf("task engine %s shutdown twice", engine.name))
	}

	log.Info(engine.ctx, "stopping")

	engine.shutdown = true
	engine.cancel()
	engine.WaitGroup.Wait()

	if engine.concurrencyTokens != nil {
		close(engine.concurrencyTokens)
	}
	log.Info(engine.ctx, "stopped")
}
