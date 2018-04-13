package kate

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"
)

//超时中间件,如果handler没在设定时间内返回，直接报告调用方超时
//timeout		定义超时的时间
//errTimeout	超时时返回的错误类型
func Timeout(timeout time.Duration, errTimeout ErrorInfo) Middleware {
	return func(h ContextHandler) ContextHandler {
		f := func(ctx context.Context, w ResponseWriter, r *Request) {
			var (
				cancel context.CancelFunc
				done   = make(chan struct{})
			)

			if timeout <= 0 {
				h.ServeHTTP(ctx, w, r)
			} else {
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()

				tw := &timeoutResponseWriter{
					ResponseWriter: w,
					h:              make(http.Header),
					errTimeout:     errTimeout,
				}

				go func() {
					h.ServeHTTP(ctx, tw, r)
					close(done)
				}()

				select {
				case <-done:
					tw.mu.Lock()
					defer tw.mu.Unlock()
					dst := w.Header()
					for k, vv := range tw.h {
						dst[k] = vv
					}
					w.WriteHeader(tw.code)
					w.Write(tw.wbuf.Bytes())
				case <-ctx.Done():
					tw.mu.Lock()
					defer tw.mu.Unlock()
					Error(ctx, w, errTimeout)
					tw.timedOut = true
				}
			}
		}
		return ContextHandlerFunc(f)
	}
}

type timeoutResponseWriter struct {
	ResponseWriter
	h          http.Header
	wbuf       bytes.Buffer
	errTimeout ErrorInfo

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

func (tw *timeoutResponseWriter) Header() http.Header {
	return tw.h
}

func (tw *timeoutResponseWriter) Write(p []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, tw.errTimeout
	}
	if !tw.wroteHeader {
		tw.writeHeader(http.StatusOK)
	}
	return tw.wbuf.Write(p)
}

func (tw *timeoutResponseWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

func (tw *timeoutResponseWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

func (tw *timeoutResponseWriter) StatusCode() int {
	return tw.code
}

func (tw *timeoutResponseWriter) RawBody() []byte {
	return tw.wbuf.Bytes()
}

func (tw *timeoutResponseWriter) WriteJSON(v interface{}) error {
	b, err := tw.EncodeJSON(v)
	if err != nil {
		return err
	}

	_, err = tw.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (tw *timeoutResponseWriter) Flush() {
	if !tw.wroteHeader {
		tw.WriteHeader(http.StatusOK)
	}
	flusher := tw.ResponseWriter.(http.Flusher)
	flusher.Flush()
}
