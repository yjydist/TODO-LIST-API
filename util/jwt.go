package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"

	"todo-list-api/model"
)

type Claims struct {
	UserID    int `json:"user_id"`
	ExpiresAt int64 `json:"exp"`
	jwt.StandardClaims
}


var jwtKey= []byte(viper.GetString("jwt.secret"))

func GenerateJWT(user model.User) (string, error) {
	return "token", nil
}

// ValidateToken 验证 JWT token 并返回声明
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 检查 token 是否过期
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
