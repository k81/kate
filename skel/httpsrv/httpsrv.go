package httpsrv

import (
	"context"
	"net/http"

	"github.com/k81/kate"
	"github.com/k81/log"
)

var (
	mctx   = context.Background()
	logger *log.Logger
)

func GetServer() *http.Server {
	logger = log.With("module", "httpsrv")
	conf := GetHttpConfig()

	// 定义中间件栈，可根据需要在下面追加
	c := kate.NewChain(
		kate.TraceId,
		kate.Logging,
		kate.Recovery,
	)

	// 注册Handler
	router := kate.NewRESTRouter(mctx)
	router.SetMaxBodyBytes(conf.MaxBodyBytes)
	router.GET("/hello", c.Then(&HelloHandler{}))

	// 生成一个http.Server对象
	server := &http.Server{
		Addr:           GetListenAddr(),
		Handler:        router,
		ReadTimeout:    conf.ReadTimeout,
		WriteTimeout:   conf.WriteTimeout,
		MaxHeaderBytes: conf.MaxHeaderBytes,
	}
	return server
}
