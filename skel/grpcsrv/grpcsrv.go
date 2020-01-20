package grpcsrv

import (
	"net"

	"github.com/cloudflare/tableflip"
	"github.com/k81/kate/log"
	"github.com/k81/kate/log/encoders/simple"
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
	logger       *zap.Logger
	accessLogger *zap.Logger
}

// Start start the grpc service
func Start(upgrader *tableflip.Upgrader, logger *zap.Logger) {
	if gService != nil {
		panic("grpcsrv start twice")
	}

	var (
		enc  = simple.NewEncoder()
		core = zapcore.NewSampler(
			log.MustNewCore(zapcore.InfoLevel, config.GRPC.LogFile, enc),
			config.GRPC.LogSampler.Tick,
			config.GRPC.LogSampler.First,
			config.GRPC.LogSampler.ThereAfter)
	)

	opts := []zap.Option{
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
	}

	accessLogger := zap.New(core, opts...)

	gService = &grpcService{
		conf:         *config.GRPC,
		upgrader:     upgrader,
		logger:       logger.Named("grpcsrv"),
		accessLogger: accessLogger,
	}

	go gService.serve()
}

// Stop stop the grpc service
func Stop() {
	if gService != nil {
		gService.stop()
	}
}

func (s *grpcService) serve() {
	var err error

	if s.listener, err = s.upgrader.Listen("tcp", s.conf.Addr); err != nil {
		s.logger.Fatal("grpc listen failed",
			zap.String("addr", s.conf.Addr),
			zap.Error(err),
		)
	}

	s.server = grpc.NewServer()

	// TODO: register grpc server impl here
	// proto.RegisterXXXServer(s.server, impl)

	s.logger.Info("grpc service started listening", zap.String("addr", s.conf.Addr))

	if err = s.server.Serve(s.listener); err != nil {
		s.logger.Fatal("failed to serve grpc service", zap.Error(err))
	}
}

func (s *grpcService) stop() {
	s.server.GracefulStop()
}
