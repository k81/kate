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
	// RingSize define the ring buffer size for timer engine
	RingSize = 3600
)

var (
	mctx = log.WithContext(context.Background(), "module", "timerengine")
)

type request struct {
	taskFunc TaskFunc
	delay    int64
	result   chan *TimerTask
}

// TimerEngine define the timer engine
// nolint:maligned
type TimerEngine struct {
	name      string
	buckets   []*list.List
	requests  chan *request
	ticker    <-chan time.Time
	tickIndex uint32
	executors *taskengine.TaskEngine
	taskIDSeq uint64
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// New create a new TimerEngine
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
		ctx:       log.WithContext(newctx, "name", name),
		cancel:    cancel,
	}

	for i := 0; i < len(e.buckets); i++ {
		e.buckets[i] = list.New()
	}

	return e
}

// Name return the name of timer engine
func (e *TimerEngine) Name() string {
	return e.name
}

// Start start the timer engine
func (e *TimerEngine) Start() {
	e.wg.Add(1)
	go e.loop()
}

// Stop stop the timer engine
func (e *TimerEngine) Stop() {
	e.cancel()
	e.wg.Wait()
}

func (e *TimerEngine) loop() {
	log.Info(e.ctx, "timer engine loop started")

	defer func() {
		if r := recover(); r != nil {
			log.Fatal(e.ctx, "panic", "error", r, "stack", utils.GetPanicStack())
		}

		e.cancel()
		e.executors.Shutdown()
		e.wg.Done()
		log.Info(e.ctx, "main loop stopped")
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

				if log.Enabled(log.LevelTrace) {
					log.Trace(e.ctx, "add timer task",
						"task_id", task.ID,
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
						tBegin = time.Now()
					)

					if log.Enabled(log.LevelTrace) {
						log.Trace(e.ctx, "schedule for tick started",
							"tick_bucket_index", tickIndex,
							"task_count", tasks.Len())
					}

					if log.Enabled(log.LevelTrace) {
						defer func() {
							log.Trace(e.ctx, "schedule for tick stopped",
								"tick_bucket_index", tickIndex,
								"duration_ms", int64(time.Since(tBegin)/time.Millisecond))
						}()
					}

					for e := tasks.Front(); e != nil; e = next {
						next = e.Next()
						// nolint:errcheck
						task := e.Value.(*TimerTask)

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

func (e *TimerEngine) nextTaskID() uint64 {
	return atomic.AddUint64(&e.taskIDSeq, 1)
}

func (e *TimerEngine) updateTickIndex() uint32 {
	tickIndex := atomic.AddUint32(&e.tickIndex, 1)
	tickIndex %= RingSize
	return tickIndex
}

func (e *TimerEngine) getTickIndex() uint32 {
	tickIndex := atomic.LoadUint32(&e.tickIndex)
	tickIndex %= RingSize
	return tickIndex
}

func (e *TimerEngine) execute(f TaskFunc) {
	e.executors.Schedule(f)
}

// Schedule schedule a timer task with delay
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
