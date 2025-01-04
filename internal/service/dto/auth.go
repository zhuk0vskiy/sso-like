package dto

type SignUpRequest struct {
	Email    string
	Password string
}

type LogInRequest struct {
	Email    string
	Password string
	Token    string
}
