package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"todo-list-api/model"
	"todo-list-api/util"
)

func UserRegister(ctx *gin.Context) {
	// user := model.User{} // Initialize a new user model
	user := model.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	user.Password = string(hashedPassword)

	// Save the user to the database
	if err := user.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user",
		})
		return
	}

	// Generate JWT token
	token, err := util.GenerateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})

}

func UserLoginIn(ctx *gin.Context) {
	user := model.User{} // Initialize a new user model
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	} // Bind the request body to the user model

	user.Login() // Login the user

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
	})

}
