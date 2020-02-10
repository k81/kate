package grpcsrv

import (
	"net"
	"path"
	"sync"

	"github.com/cloudflare/tableflip"
	"github.com/k81/kate/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"__PACKAGE_NAME__/config"
)

var gService *grpcService

type grpcService struct {
	conf         config.GRPCConfig
	upgrader     *tableflip.Upgrader
	listener     net.Listener
	server       *grpc.Server
	wg           sync.WaitGroup
	logger       *zap.Logger
	accessLogger *zap.Logger
}

// Start start the grpc service
func Start(upgrader *tableflip.Upgrader, logger *zap.Logger) {
	if gService != nil {
		panic("grpcsrv start twice")
	}

	gService = &grpcService{
		conf:     *config.GRPC,
		upgrader: upgrader,
		logger:   logger.Named("grpcsrv"),
	}
	gService.start()
}

// Stop stop the grpc service
func Stop() {
	if gService != nil {
		gService.stop()
	}
}

func (s *grpService) start() {
	var (
		enc  = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		core = log.MustNewCore(zapcore.InfoLevel, path.Join(config.Main.LogDir, s.conf.LogFile), enc)
		err  error
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

	if s.listener, err = s.upgrader.Listen("tcp", s.conf.Addr); err != nil {
		s.logger.Fatal("grpc listen failed",
			zap.String("addr", s.conf.Addr),
			zap.Error(err),
		)
	}

	s.server = grpc.NewServer()

	// TODO: register grpc server impl here
	// proto.RegisterXXXServer(s.server, impl)

	s.wg.Add(1)
	go s.serve()
}

func (s *grpcService) serve() {
	defer func() {
		s.wg.Done()
		s.logger.Info("grpc service stopped")
	}()

	s.logger.Info("grpc service started listening", zap.String("addr", s.conf.Addr))

	if err := s.server.Serve(s.listener); err != nil {
		s.logger.Fatal("failed to serve grpc service", zap.Error(err))
	}
}

func (s *grpcService) stop() {
	s.server.GracefulStop()
	s.wg.Wait()
}
