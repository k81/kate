package httpsrv

import (
	"context"
	"net/http"

	"github.com/k81/kate"
	"go.uber.org/zap"

	"__PACKAGE_NAME__/config"
)

var logger *zap.Logger

// ListenAndServe start the http server and wait for exit
func ListenAndServe(l *zap.Logger) {
	logger = l.With(zap.String("module", "httpsrv"))
	// 定义中间件栈，可根据需要在下面追加
	c := kate.NewChain(
		Logging,
		Recovery,
	)

	// 注册Handler
	router := kate.NewRESTRouter(context.Background(), logger)
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

	logger.Info("http service started", zap.String("listen_addr", config.HTTP.Addr))

	if err := server.ListenAndServe(); err != nil {
		logger.Error("serve error", zap.Error(err))
	}
}
