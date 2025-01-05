package token

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type TokenInterface interface {
	Create(username string, duration time.Duration) (string, error)
	Validate(token string) (*Payload, error)
}

type Payload struct {
	Id        uuid.UUID
	Username  string
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		Id:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {

	if time.Now().After(payload.ExpiredAt) {
		return errors.Errorf("token has expired")
	}
	return nil
}
