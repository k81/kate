package httpsrv

import (
	"context"
	"net"
	"net/http"

	"github.com/cloudflare/tableflip"
	"github.com/k81/kate"
	"github.com/k81/kate/log"
	"github.com/k81/kate/log/encoders/simple"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"__PACKAGE_NAME__/config"
)

var gService *httpService

type httpService struct {
	conf         config.HTTPConfig
	upgrader     *tableflip.Upgrader
	listener     net.Listener
	server       *http.Server
	logger       *zap.Logger
	accessLogger *zap.Logger
}

// Start start the http service
func Start(upgrader *tableflip.Upgrader, logger *zap.Logger) {
	if gService != nil {
		panic("httpsrv start twice")
	}

	var (
		enc  = simple.NewEncoder()
		core = zapcore.NewSampler(
			log.MustNewCore(zapcore.InfoLevel, config.HTTP.LogFile, enc),
			config.HTTP.LogSampler.Tick,
			config.HTTP.LogSampler.First,
			config.HTTP.LogSampler.ThereAfter)
	)

	opts := []zap.Option{
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
	}

	accessLogger := zap.New(core, opts...)

	gService = &httpService{
		conf:         *config.HTTP,
		upgrader:     upgrader,
		logger:       logger.Named("httpsrv"),
		accessLogger: accessLogger,
	}

	go gService.serve()
}

// Stop stop the http service
func Stop() {
	if gService != nil {
		gService.stop()
	}
}

func (s *httpService) serve() {
	var err error

	// 定义中间件栈，可根据需要在下面追加
	c := kate.NewChain(
		Logging,
		Recovery,
	)

	// 注册Handler
	router := kate.NewRESTRouter(context.Background(), s.logger)
	router.SetMaxBodyBytes(s.conf.MaxBodyBytes)
	router.GET("/hello", c.Then(&HelloHandler{}))

	// 生成一个http.Server对象
	s.server = &http.Server{
		Addr:           s.conf.Addr,
		Handler:        router,
		ReadTimeout:    s.conf.ReadTimeout,
		WriteTimeout:   s.conf.WriteTimeout,
		MaxHeaderBytes: s.conf.MaxHeaderBytes,
	}

	if s.listener, err = s.upgrader.Listen("tcp", s.conf.Addr); err != nil {
		s.logger.Fatal("http listen failed",
			zap.String("addr", s.conf.Addr),
			zap.Error(err),
		)
	}

	s.logger.Info("http service started listening", zap.String("addr", s.conf.Addr))

	if err = s.server.Serve(s.listener); err != nil {
		s.logger.Fatal("failed to serve http service", zap.Error(err))
	}
}

func (s *httpService) stop() {
	if err := s.server.Shutdown(context.TODO()); err != nil {
		s.logger.Error("http service shutdown failed", zap.Error(err))
	}
}
