package httpclient

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"context"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"
)

var (
	defaultTransport = &http.Transport{
		MaxIdleConnsPerHost: 16,
	}

	DefaultTimeout = 15 * time.Second
	Debug          = true
)

type ErrStatusCode struct {
	StatusCode int
	Status     string
}

func (e ErrStatusCode) Error() string {
	return fmt.Sprintf("Error: %v", e.Status)
}

type RetryClassifier struct{}

func (c RetryClassifier) Classify(err error) retrier.Action {
	if err == nil {
		return retrier.Succeed
	}

	if ne, ok := err.(net.Error); ok && ne.Temporary() {
		return retrier.Retry
	}

	if he, ok := err.(*ErrStatusCode); ok && he.StatusCode >= 500 {
		return retrier.Retry
	}
	return retrier.Fail
}

func disableRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func DoWithTransport(ctx context.Context, method, url, contentType, reqBody string, timeout time.Duration, followRedirect bool, transport http.RoundTripper) (respBody string, err error) {
	var (
		request  *http.Request
		response *http.Response
		respData []byte
		client   = http.Client{
			Timeout: timeout,
		}
	)

	if transport == nil {
		transport = defaultTransport
	}

	client.Transport = transport

	if request, err = http.NewRequest(method, url, strings.NewReader(reqBody)); err != nil {
		log.Error(ctx, "create http request", "method", method, "url", url, "request_body", reqBody, "error", err)
		return
	}

	traceId := utils.GetString(ctx.Value("X-Trace-Id"))
	if traceId == "" {
		traceId = utils.FastUUIDStr()
	}
	request.Header.Add("X-Trace-Id", traceId)

	if contentType != "" {
		request.Header.Add("Content-Type", contentType)
	}

	if !followRedirect {
		client.CheckRedirect = disableRedirect
	}

	if Debug {
		log.Debug(ctx, "http request begin", "trace_id", traceId, "method", method, "url", url, "body", reqBody)
	}

	tBegin := time.Now()
	if response, err = client.Do(request); err != nil {
		log.Error(ctx, "do http request", "method", method, "url", url, "request_body", reqBody, "error", err)
		return
	}
	defer response.Body.Close()

	if respData, err = ioutil.ReadAll(response.Body); err != nil {
		log.Error(ctx, "read http response body", "method", method, "url", url, "request_body", reqBody, "error", err)
		return
	}

	respBody = string(respData)

	if Debug {
		log.Debug(ctx, "http request end", "trace_id", traceId, "status_code", response.StatusCode, "duration_ms", time.Since(tBegin)/time.Millisecond)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Error(ctx, "bad http status code",
			"method", method,
			"url", url,
			"status", response.StatusCode,
			"request_body", reqBody,
			"response_body", respBody)

		err = &ErrStatusCode{response.StatusCode, response.Status}
		return
	}

	return
}

func Do(ctx context.Context, method, url, contentType, reqBody string, timeout time.Duration, followRedirect bool) (respBody string, err error) {
	return DoWithTransport(ctx, method, url, contentType, reqBody, timeout, followRedirect, nil)
}

func TryCallWithTransport(ctx context.Context, method, url, contentType, reqBody string, timeout time.Duration, backoff []time.Duration, transport http.RoundTripper) (respBody string, err error) {
	r := retrier.New(backoff, RetryClassifier{})

	err = r.Run(func() error {
		if respBody, err = DoWithTransport(ctx, method, url, contentType, reqBody, timeout, false, transport); err != nil {
			return err
		}

		return nil
	})

	return
}

func TryCall(ctx context.Context, method, url, contentType, reqBody string, timeout time.Duration, backoff []time.Duration) (respBody string, err error) {
	return TryCallWithTransport(ctx, method, url, contentType, reqBody, timeout, backoff, nil)
}
