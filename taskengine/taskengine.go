package taskengine

import (
	"fmt"
	"sync"

	"context"

	"github.com/k81/kate/utils"
	"github.com/k81/log"
)

type TaskEngine struct {
	sync.WaitGroup
	name              string
	shutdown          bool
	concurrencyTokens chan struct{}
	logger            *log.Logger
	ctx               context.Context
	cancel            context.CancelFunc
}

func New(ctx context.Context, name string, concurrencyLevel int) *TaskEngine {
	var (
		newctx, cancel = context.WithCancel(ctx)
		logger         = log.With("taskengine", name)
	)

	engine := &TaskEngine{
		name:   name,
		logger: logger,
		ctx:    newctx,
		cancel: cancel,
	}

	if concurrencyLevel > 0 {
		engine.concurrencyTokens = make(chan struct{}, concurrencyLevel)
	}
	return engine
}

func (engine *TaskEngine) Schedule(task Task) bool {
	if engine.shutdown {
		engine.logger.Error(engine.ctx, "already stopped, should not schedule new task")
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
			engine.logger.Error(engine.ctx, "task panic:", "error", r, "stack", utils.GetPanicStack())
		}
		if engine.concurrencyTokens != nil {
			<-engine.concurrencyTokens
		}
		engine.Done()
	}()

	task.Run()
}

func (engine *TaskEngine) Shutdown() {
	if engine.shutdown {
		panic(fmt.Sprintf("task engine %s shutdown twice", engine.name))
	}

	engine.logger.Info(engine.ctx, "stopping")

	engine.shutdown = true
	engine.cancel()
	engine.WaitGroup.Wait()

	if engine.concurrencyTokens != nil {
		close(engine.concurrencyTokens)
	}
	engine.logger.Info(engine.ctx, "stopped")
}
