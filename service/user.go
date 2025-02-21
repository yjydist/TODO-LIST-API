package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"todo-list-api/model"
	"todo-list-api/util"
)

func UserRegister(ctx *gin.Context) {
	// 将请求的 JSON 数据绑定到 User 结构体
	user := model.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 将 User 结构体中的密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	user.Password = string(hashedPassword)

	// 保存用户
	if err := user.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user",
		})
		return
	}

	 // 生成 JWT token
	token, err := util.GenerateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// 返回成功信息
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})

}

func UserLoginIn(ctx *gin.Context) {
	// 将请求的 JSON 数据绑定到 User 结构体
	user := model.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 查找用户
	storedUser, err := model.FindUserByUsername(user.Name)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// 生成 JWT token
	token, err := util.GenerateJWT(storedUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// 返回成功信息
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"token":   token,
	})
}
