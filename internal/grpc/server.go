package grpc

import (
	"fmt"
	"net"
)

// func (a *App) MustRun() {
// 	if err := a.Run(); err != nil {
// 		panic(err)
// 	}
// }

// Run runs gRPC server.
func Run(a *GrpcApp) error {
	const op = "grpcapp.Run"

	// Создаём listener, который будет слушить TCP-сообщения, адресованные
	// Нашему gRPC-серверу
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Infof("grpc server started", l.Addr().String())
	fmt.Println("grpc server started", l.Addr().String())

	// Запускаем обработчик gRPC-сообщений
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
