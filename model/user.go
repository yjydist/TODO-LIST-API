package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);unique;not null"`
}

// TableName 指定表名
func (u *User) TableName() string {
	return "users"
}

// Save 保存用户到数据库
func (u *User) Save() error {
	return DB.Create(u).Error
}

// FindUserByEmail 通过邮箱查找用户
func FindUserByEmail(email string) (User, error) {
	var user User
	result := DB.Where("email = ?", email).First(&user)
	return user, result.Error
}

// FindUserByUsername 通过用户名查找用户
func FindUserByUsername(username string) (User, error) {
	var user User
	result := DB.Where("name = ?", username).First(&user)
	return user, result.Error
}

// 用于JWT的Claims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(user User) (string, error) {
	// 设置JWT声明
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// 创建JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// HashPassword 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash 验证密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
