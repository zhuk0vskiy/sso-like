package dto

type SignUpRequest struct {
	Email    string
	Password string
}

type LogInRequest struct {
	Email    string
	Password string
	AppId    int64
}
