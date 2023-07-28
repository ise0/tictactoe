package authjwt

import (
	"api/src/shared/logger"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

type claims struct {
	jwt.StandardClaims
	UserId int `json:"uid"`
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func SignAuthJwt(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{UserId: userId})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		logger.L.Errorw("Failed to sign string")
		return "", err
	}
	return tokenString, nil
}

func ParseAuthJwt(t string) (int, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})
	if err != nil {
		logger.L.Infow("Failed to parse jwt", "error", err)
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims["uid"].(float64); ok {
			return int(userId), nil
		}
	}
	logger.L.Infow("Failed to parse jwt", "error", err)
	return 0, fmt.Errorf("failed to parse jwt")
}
