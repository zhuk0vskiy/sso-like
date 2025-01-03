package dto

type InsertUserRequest struct {
	Email    string
	PassHash []byte
}

type GetUserRequest struct {
	Email string
}
