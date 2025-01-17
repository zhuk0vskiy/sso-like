package handler

import (
	"context"
	// authGrpc "sso-like/internal/grpc/auth"
	authService "sso-like/internal/service/auth"
	dtoService "sso-like/internal/service/dto"

	// "sso-like/internal/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ssov1 "github.com/zhuk0vskiy/protos/gen/go/sso"
)

type ServerApi struct {
	ssov1.UnimplementedAuthServer
	Auth authService.AuthInterface
}

func (s *ServerApi) LogIn(ctx context.Context, in *ssov1.LogInRequest) (*ssov1.LogInResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	token, err := s.Auth.LogIn(ctx, &dtoService.LogInRequest{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
		Token:    in.GetToken(),
	})

	if err != nil {

		// if errors.Is(err, authGrpc.ErrInvalidCredentials) {
		// 	return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		// }

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LogInResponse{Token: token}, nil
}

func (s *ServerApi) SignUp(ctx context.Context, in *ssov1.SignUpRequest) (*ssov1.SignUpResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	totpSecret, err := s.Auth.SignUp(ctx, &dtoService.SignUpRequest{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	})
	if err != nil {

		// if errors.Is(err, storage.ErrUserExists) {
		//     return nil, status.Error(codes.AlreadyExists, "user already exists")
		// }

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.SignUpResponse{
		TotpSecret: totpSecret,
		TotpUrl:    "",
	}, nil
}
