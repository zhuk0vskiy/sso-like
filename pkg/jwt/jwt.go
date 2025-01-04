package jwt

import (
	"sso-like/internal/model"
	"time"

	"github.com/golang-jwt/jwt"
)

var JWTKey = "sMMkd8fLdv"

func NewToken(user *model.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	// Подписываем токен, используя секретный ключ приложения
	tokenString, err := token.SignedString([]byte(JWTKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
