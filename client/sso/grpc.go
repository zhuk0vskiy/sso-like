package main

import (
	"context"
	"fmt"

	ssov1 "github.com/zhuk0vskiy/protos/gen/go/sso"
	"google.golang.org/grpc"
)

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

	// req := &ssov1.SignUpRequest{
	// 	Email:    "a",
	// 	Password: "1",
	// }

	// response, err := grpcClient.SignUp(context.Background(), req)

	// if err != nil {
	// 	fmt.Printf("fail to dial: %v", err)
	// 	return
	// }

	// fmt.Println(response.TotpSecret)

	reqq := &ssov1.LogInRequest{
		Email:    "a",
		Password: "1",
		Token:    "743590",
	}
	ress, err := grpcClient.LogIn(context.Background(), reqq)
	if err != nil {
		fmt.Printf("fail to dial: %v", err)
		return
	}

	fmt.Println(ress.Token)
}
