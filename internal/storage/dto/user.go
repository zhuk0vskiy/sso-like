package dto

type InsertUserRequest struct {
	Email      string
	Password   []byte
	TotpSecret []byte
}

type GetUserRequest struct {
	Email string
}
