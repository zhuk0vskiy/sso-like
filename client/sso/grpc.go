package main

import (
	"context"
	"fmt"

	ssov1 "github.com/zhuk0vskiy/protos/gen/go/sso"
	"google.golang.org/grpc"
)

// type Client struct {
// 	api ssov1.AuthClient
// 	log logger.Logger
// }

// func NewClient(
// 	ctx context.Context,
// 	log logger.Logger,
// 	addr string, // Адрес SSO-сервера
// 	timeout time.Duration, // Таймаут на выполнение каждой попытки
// 	retriesCount int, // Количетсво повторов
// ) (*Client, error) {
// 	const op = "grpc.New"

// 	// Опции для интерсептора grpcretry
// 	retryOpts := []grpcretry.CallOption{
// 		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
// 		grpcretry.WithMax(uint(retriesCount)),
// 		grpcretry.WithPerRetryTimeout(timeout),
// 	}

// 	// Опции для интерсептора grpclog
// 	logOpts := []grpclog.Option{
// 		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
// 	}

// 	// Создаём соединение с gRPC-сервером SSO для клиента
// 	cc, err := grpc.DialContext(ctx, addr,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithChainUnaryInterceptor(
// 			grpclog.UnaryClientInterceptor(logger.InterceptorLogger(log), logOpts...),
// 			grpcretry.UnaryClientInterceptor(retryOpts...),
// 		))
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	// Создаём gRPC-клиент SSO/Auth

// 	return &Client{
// 		api: grpcClient,
// 	}, nil
// }

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	// args := os.Args
	conn, err := grpc.Dial("127.0.0.1:44044", opts...)

	if err != nil {
		// grpclog.Fatal
		fmt.Printf("fail to dial: %v", err)
		return
	}

	defer conn.Close()

	grpcClient := ssov1.NewAuthClient(conn)

	// client := pb.NewReverseClient(conn)
	req := &ssov1.SignUpRequest{
		Email:    "1",
		Password: "1",
	}
	// request := &pb.Request{
	//     Message: args[1],
	// }
	response, err := grpcClient.SignUp(context.Background(), req)
	// response, err := client.Do(context.Background(), request)

	if err != nil {
		fmt.Printf("fail to dial: %v", err)
		return
	}

	fmt.Println(response.UserId)

	// reqq := &ssov1.LogInRequest{
	// 	Email:    "1",
	// 	Password: "1",
	// 	AppId: 1,
	// }
	// request := &pb.Request{
	//     Message: args[1],
	// }
	// response, err = grpcClient.LogIn(context.Background(), reqq)
	// // response, err := client.Do(context.Background(), request)

	// if err != nil {
	// 	fmt.Printf("fail to dial: %v", err)
	// 	return
	// }

	// fmt.Println(response.UserId)
}
