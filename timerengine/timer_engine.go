package timerengine

import (
	"container/list"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/k81/kate/taskengine"
	"github.com/k81/kate/utils"
	"github.com/k81/log"
)

const (
	RingSize = 3600
)

var (
	mctx = context.Background()
)

type request struct {
	taskFunc TaskFunc
	delay    int64
	result   chan *TimerTask
}

type TimerEngine struct {
	name      string
	buckets   []*list.List
	requests  chan *request
	ticker    <-chan time.Time
	tickIndex uint32
	executors *taskengine.TaskEngine
	taskIdSeq uint64
	logger    *log.Logger
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func New(name string, concurrencyLevel int) *TimerEngine {
	var (
		newctx, cancel = context.WithCancel(mctx)
	)

	e := &TimerEngine{
		name:      name,
		ticker:    time.Tick(time.Second),
		buckets:   make([]*list.List, RingSize),
		requests:  make(chan *request, 1024),
		executors: taskengine.New(newctx, name, concurrencyLevel),
		logger:    log.With("name", name),
		ctx:       newctx,
		cancel:    cancel,
	}

	for i := 0; i < len(e.buckets); i++ {
		e.buckets[i] = list.New()
	}

	return e
}

func (e *TimerEngine) Name() string {
	return e.name
}

func (e *TimerEngine) Start() {
	e.wg.Add(1)
	go e.loop()
}

func (e *TimerEngine) Stop() {
	e.cancel()
	e.wg.Wait()
}

func (e *TimerEngine) loop() {
	e.logger.Info(e.ctx, "timer engine loop started")

	defer func() {
		if r := recover(); r != nil {
			e.logger.Fatal(e.ctx, "panic", "error", r, "stack", utils.GetPanicStack())
		}

		e.cancel()
		e.executors.Shutdown()
		e.wg.Done()
		e.logger.Info(e.ctx, "main loop stopped")
	}()

	for {
		select {
		case <-e.ctx.Done():
			return
		case req := <-e.requests:
			{
				var (
					tickIndex   = e.getTickIndex()
					offset      = int(req.delay) + int(tickIndex)
					cycleNum    = offset / RingSize
					bucketIndex = offset % RingSize
					bucket      = e.buckets[bucketIndex]
					task        = newTimerTask(e, cycleNum, req.taskFunc)
				)
				bucket.PushBack(task)

				req.result <- task

				if e.logger.Enabled(log.LevelTrace) {
					e.logger.Trace(e.ctx, "add timer task",
						"task_id", task.Id,
						"delay", req.delay,
						"tick_bucket_index", tickIndex,
						"task_bucket_index", bucketIndex,
						"cycleNum", cycleNum,
					)
				}
			}
		case <-e.ticker:
			{
				var (
					tickIndex = e.updateTickIndex()
					tasks     = e.buckets[tickIndex]
				)

				go func(tasks *list.List, tickIndex uint32) {
					var (
						next   *list.Element
						task   *TimerTask
						tBegin = time.Now()
					)

					if e.logger.Enabled(log.LevelTrace) {
						e.logger.Trace(e.ctx, "schedule for tick started", "tick_bucket_index", tickIndex, "task_count", tasks.Len())
					}

					if e.logger.Enabled(log.LevelTrace) {
						defer func() {
							e.logger.Trace(e.ctx, "schedule for tick stopped", "tick_bucket_index", tickIndex, "duration_ms", int64(time.Since(tBegin)/time.Millisecond))
						}()
					}

					for e := tasks.Front(); e != nil; e = next {
						next = e.Next()
						task = e.Value.(*TimerTask)

						if task.ready() {
							task.dispose()
							tasks.Remove(e)
						}
					}
				}(tasks, tickIndex)
			}
		}
	}
}

func (e *TimerEngine) nextTaskId() uint64 {
	return atomic.AddUint64(&e.taskIdSeq, 1)
}

func (e *TimerEngine) updateTickIndex() uint32 {
	tickIndex := atomic.AddUint32(&e.tickIndex, 1)
	tickIndex = tickIndex % RingSize
	return tickIndex
}

func (e *TimerEngine) getTickIndex() uint32 {
	tickIndex := atomic.LoadUint32(&e.tickIndex)
	tickIndex = tickIndex % RingSize
	return tickIndex
}

func (e *TimerEngine) execute(f TaskFunc) {
	e.executors.Schedule(f)
}

func (e *TimerEngine) Schedule(f TaskFunc, delay int64) (task *TimerTask) {
	if delay <= 0 {
		task = newTimerTask(e, 0, f)
		task.dispose()
		return
	}

	req := &request{
		taskFunc: f,
		delay:    delay,
		result:   make(chan *TimerTask, 1),
	}

	select {
	case <-e.ctx.Done():
		return
	case e.requests <- req:
	}

	select {
	case <-e.ctx.Done():
	case task = <-req.result:
	}
	return
}
