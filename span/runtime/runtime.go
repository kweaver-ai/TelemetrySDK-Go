package runtime

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/open_standard"
)

var (
	defaultInternal = 10 * time.Second
	defaultMaxLog   = 40
)

// Runtime read data from channel and write data in a single goroutine
type Runtime struct {
	cache   chan field.LogSpan
	builder func(func(field.LogSpan), context.Context) field.LogSpan
	stop    chan int
	wg      *sync.WaitGroup
	// id      uint64
	close     bool
	closeLock sync.RWMutex
	runLock   sync.Mutex
	w         open_standard.Writer
	once      sync.Once
	logs      []field.LogSpan
	size      int
	maxLog    int
	internal  time.Duration
}

// NewRuntime return a runtime
func NewRuntime(w open_standard.Writer, builder func(func(field.LogSpan), context.Context) field.LogSpan) *Runtime {
	r := &Runtime{
		cache:   make(chan field.LogSpan, 100),
		builder: builder,
		stop:    make(chan int, 1),
		wg:      &sync.WaitGroup{},
		// false = running, true = closed or closing
		close: false,
		// protect the close value
		closeLock: sync.RWMutex{},
		// represent the state of runtime thread, runtime thread will lock this until thread over
		runLock:  sync.Mutex{},
		w:        w,
		once:     sync.Once{},
		size:     0,
		logs:     make([]field.LogSpan, 0, defaultMaxLog+1),
		maxLog:   defaultMaxLog,
		internal: defaultInternal,
	}

	return r
}

// Children return a logger span
// if Runtime has been close return nil
// user should return span's onwership after Span is useless by Span.Signal()
func (r *Runtime) Children(ctx context.Context) field.LogSpan {
	//remove read lock
	r.closeLock.RLock()
	defer r.closeLock.RUnlock()

	if r.close {
		return nil
	}
	s := r.builder(r.transfer, ctx)
	r.wg.Add(1)
	return s
}

// Signal stop runtime thread
func (r *Runtime) Signal() {
	r.closeLock.Lock()
	r.close = true
	r.wg.Wait()
	if r.size > 0 {
		r.forceWrite()
	}
	r.once.Do(func() {
		r.stop <- 0
		close(r.cache)
		_ = r.w.Close()
	})
	r.runLock.Lock()
	//nolint 这里不去掉是因为git提交记录[bugfix](runtime) 关闭runtime线程与runtime运行线程同步问题可能导致数据丢失
	r.runLock.Unlock()
	r.closeLock.Unlock()
}

func (r *Runtime) transfer(s field.LogSpan) {
	r.cache <- s
}

// Run will deal Runtime's span in current go runtine
// Run will block the go runtine, so you may need to use a go runtine to run
// Runtime will close when all span has been did when receive system signal or do Signal()
// Maybe runtime should not close by signal
func (r *Runtime) Run() {
	r.runLock.Lock()
	defer r.runLock.Unlock()
	Ticker := time.NewTicker(r.internal)
	for {
		select {
		case s, ok := <-(r.cache):
			// 关闭之后退出循环。
			if !ok {
				// 发完最后的数据关闭Exporter。
				err := r.w.Close()
				if err != nil {
					log.Println(field.GenerateSpecificError(err))
				}
				return
			}
			r.logs = append(r.logs, s)
			r.size++
			r.wg.Done()
			// 超过上限发送。
			if r.size >= r.maxLog {
				r.forceWrite()
			}
		// 定时发送。
		case <-Ticker.C:
			if r.size > 0 {
				r.forceWrite()
			}
		}
	}
}

func (r *Runtime) forceWrite() {
	err := r.w.Write(r.logs)
	if err != nil {
		log.Println(field.GenerateSpecificError(err))
	}
	// 发送完之后清空队列。
	for _, v := range r.logs {
		v.Free()
	}
	r.size = 0
	r.logs = make([]field.LogSpan, 0, r.maxLog+1)
}

func (r *Runtime) SetUploadInternalAndMaxLog(Internal time.Duration, MaxLog int) {
	r.internal = Internal
	r.maxLog = MaxLog
	r.logs = make([]field.LogSpan, 0, r.maxLog+1)
}
