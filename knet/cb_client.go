package knet

import (
	"context"
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreakerClient struct {
	breaker *gobreaker.CircuitBreaker
	Client
}

func NewCircuitBreakerClient(client Client, breaker *gobreaker.CircuitBreaker) *CircuitBreakerClient {
	return &CircuitBreakerClient{
		Client:  client,
		breaker: breaker,
	}
}

func (c *CircuitBreakerClient) Dial(addr string) (err error) {
	_, err = c.breaker.Execute(func() (interface{}, error) { return nil, c.Client.Dial(addr) })
	return
}

func (c *CircuitBreakerClient) Call(ctx context.Context, req Request) (resp Response, err error) {
	var reply interface{}

	reply, err = c.breaker.Execute(func() (interface{}, error) { return c.Client.Call(ctx, req) })
	if err == nil {
		resp = reply.(Response)
	}
	return
}

func (c *CircuitBreakerClient) CallWithTimeout(ctx context.Context, req Request, timeout time.Duration) (resp Response, err error) {
	var reply interface{}

	reply, err = c.breaker.Execute(func() (interface{}, error) { return c.Client.CallWithTimeout(ctx, req, timeout) })
	if err == nil {
		resp = reply.(Response)
	}
	return
}

func (c *CircuitBreakerClient) Send(ctx context.Context, msg Message) (err error) {
	_, err = c.breaker.Execute(func() (interface{}, error) { return nil, c.Client.Send(ctx, msg) })
	return
}

func (c *CircuitBreakerClient) SendWithTimeout(ctx context.Context, msg Message, timeout time.Duration) (err error) {
	_, err = c.breaker.Execute(func() (interface{}, error) { return nil, c.Client.SendWithTimeout(ctx, msg, timeout) })
	return
}
