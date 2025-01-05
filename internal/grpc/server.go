package grpc

import (
	"fmt"
	"net"
)


func Run(a *GrpcApp) error {
	const op = "grpcapp.Run"


	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Infof("grpc server started", l.Addr().String())
	fmt.Println("grpc server started", l.Addr().String())


	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
