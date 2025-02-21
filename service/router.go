package service

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouter(server *gin.Engine) {
	server.POST("/register", UserRegister)
	server.POST("/login", UserLoginIn)

	// authServer := server.Group("/")

}
