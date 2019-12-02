package taskengine

import (
	"fmt"
	"sync"

	"context"

	"go.uber.org/zap"
)

// TaskEngine define the task engine
type TaskEngine struct {
	name              string
	concurrencyTokens chan struct{}
	ctx               context.Context
	cancel            context.CancelFunc
	logger            *zap.Logger
	shutdown          bool
	sync.WaitGroup
}

// New create a new task engine
func New(ctx context.Context, name string, concurrencyLevel int, logger *zap.Logger) *TaskEngine {
	newctx, cancel := context.WithCancel(ctx)

	engine := &TaskEngine{
		name:   name,
		ctx:    newctx,
		cancel: cancel,
		logger: logger.With(zap.String("taskengine", name)),
	}

	if concurrencyLevel > 0 {
		engine.concurrencyTokens = make(chan struct{}, concurrencyLevel)
	}
	return engine
}

// Schedule schedule a task running on engine
func (engine *TaskEngine) Schedule(task Task) bool {
	if engine.shutdown {
		engine.logger.Error("already stopped, should not schedule new task")
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
			engine.logger.Error("task panic:",
				zap.Any("error", r),
				zap.Stack("stack"),
			)
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

	engine.logger.Info("stopping")

	engine.shutdown = true
	engine.cancel()
	engine.WaitGroup.Wait()

	if engine.concurrencyTokens != nil {
		close(engine.concurrencyTokens)
	}
	engine.logger.Info("stopped")
}
