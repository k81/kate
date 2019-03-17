package httpsrv

import (
	"context"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/k81/kate"
	"github.com/k81/log"

	"__PROJECT_DIR__/config"
)

var (
	mctx = log.WithContext(context.Background(), "module", "httpsrv")
)

// ListenAndServe start the http server and wait for exit
func ListenAndServe() {
	// 定义中间件栈，可根据需要在下面追加
	c := kate.NewChain(
		kate.Logging,
		kate.Recovery,
	)

	// 注册Handler
	router := kate.NewRESTRouter(mctx)
	router.SetMaxBodyBytes(config.HTTP.MaxBodyBytes)
	router.GET("/hello", c.Then(&HelloHandler{}))

	// 生成一个http.Server对象
	server := &http.Server{
		Addr:           config.HTTP.Addr,
		Handler:        router,
		ReadTimeout:    config.HTTP.ReadTimeout,
		WriteTimeout:   config.HTTP.WriteTimeout,
		MaxHeaderBytes: config.HTTP.MaxHeaderBytes,
	}

	log.Info(mctx, "http service started", "listen_addr", config.HTTP.Addr)

	if err := gracehttp.Serve(server); err != nil {
		log.Error(mctx, "serve error", "error", err)
	}
}
