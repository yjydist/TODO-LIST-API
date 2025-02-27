package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// Claims JWT Claims结构
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Rate limiting functionality moved to rate_limit.go

// GenerateJWT 生成JWT令牌
func GenerateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用统一的JWT密钥配置
	jwtKey := []byte(viper.GetString("jwt.secret"))
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWT 验证JWT令牌
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// 使用统一的JWT密钥配置
	jwtKey := []byte(viper.GetString("jwt.secret"))

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// AuthMiddleware 身份验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		// 检查token格式
		parts := strings.Split(authHeader, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token format",
			})
			c.Abort()
			return
		}

		// 验证token
		tokenString := parts[1]
		claims, err := ValidateJWT(tokenString) // 使用本地函数而不是model包中的函数
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// 将用户ID保存到上下文
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
