package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go-grpc-sso/internal/domain/models"
	"time"
)

// NewToken creates new JWT token for given user and app
func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	//todo test
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Parse parses given token and return its claims
func Parse(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse token")
	}

	return claims, nil
}
