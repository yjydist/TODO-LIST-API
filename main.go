package main

import (
	"github.com/gin-gonic/gin"

	"todo-list-api/service"
)

func main() {
	r := gin.Default()
	//
	r.POST("/register", service.UserRegister)
	r.POST("/login", service.UserLoginIn)
	r.Run(":8080")

}
