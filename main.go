package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"todo-list-api/config"
	"todo-list-api/middleware"
	"todo-list-api/model"
	"todo-list-api/service"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	model.InitDB()

	// 创建 Gin 实例
	r := gin.Default()

	// 应用限速中间件
	r.Use(middleware.RateLimitMiddleware())

	// 注册路由
	// 用户注册和登录
	r.POST("/register", service.UserRegister)
	r.POST("/login", service.UserLoginIn)

	// 需要认证的路由
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// 待办事项管理
		auth.POST("/todos", service.TodoCreate)
		auth.PUT("/todos/:id", service.TodoUpdate)
		auth.DELETE("/todos/:id", service.TodoDelete)
		auth.GET("/todos", service.TodoGet)
		auth.GET("/todos/:id", service.TodoGetByID)
		auth.PATCH("/todos/:id/complete", service.TodoComplete)
	}

	// 启动服务器
	port := fmt.Sprintf(":%s", "8080")
	r.Run(port)
}
