package model

type User struct {
	Id       int64
	Email    string
	PassHash []byte
}
