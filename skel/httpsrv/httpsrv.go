package httpsrv

import (
	"context"
	"net"
	"net/http"
	"path"
	"sync"

	"github.com/cloudflare/tableflip"
	"github.com/k81/kate"
	"github.com/k81/kate/log"
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
	wg           sync.WaitGroup
	logger       *zap.Logger
	accessLogger *zap.Logger
}

// Start start the http service
func Start(upgrader *tableflip.Upgrader, logger *zap.Logger) {
	if gService != nil {
		panic("httpsrv start twice")
	}

	gService = &httpService{
		conf:     *config.HTTP,
		upgrader: upgrader,
		logger:   logger.Named("httpsrv"),
	}
	gService.start()
}

// Stop stop the http service
func Stop() {
	if gService != nil {
		gService.stop()
	}
}

func (s *httpService) start() {
	var (
		enc  = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		core = log.MustNewCore(zapcore.InfoLevel, path.Join(config.Main.LogDir, s.conf.LogFile), enc)
	)

	if s.conf.LogSampler.Enabled {
		core = zapcore.NewSampler(
			core,
			s.conf.LogSampler.Tick,
			s.conf.LogSampler.First,
			s.conf.LogSampler.ThereAfter,
		)
	}

	opts := []zap.Option{
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
	}

	s.accessLogger = zap.New(core, opts...)

	s.wg.Add(1)
	go s.serve()
}

func (s *httpService) serve() {
	defer func() {
		s.wg.Done()
		s.logger.Info("http service stopped")
	}()

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
	s.wg.Wait()
}
