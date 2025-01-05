package grpc

import (
	"sso-like/internal/grpc/handler"
	authService "sso-like/internal/service/auth"
	"sso-like/pkg/logger"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	// "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	ssov1 "github.com/zhuk0vskiy/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

func Register(gRPCServer *grpc.Server, auth authService.AuthInterface) {
	ssov1.RegisterAuthServer(gRPCServer, &handler.ServerApi{Auth: auth})
}

type GrpcApp struct {
	log        logger.Interface
	gRPCServer *grpc.Server
	port       int
}

func NewGrpcApp(log logger.Interface, authIntf authService.AuthInterface, port int) *GrpcApp {
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			// Логируем информацию о панике с уровнем Error
			log.Errorf("Recovered from panic:", err)

			// Можете либо честно вернуть клиенту содержимое паники
			// Либо ответить - "internal error", если не хотим делиться внутренностями
			return status.Errorf(codes.Internal, "internal error", err)
		}),
	}

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(logger.InterceptorLogger(log), loggingOpts...),
	))

	Register(gRPCServer, authIntf)

	return &GrpcApp{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
