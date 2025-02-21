package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"todo-list-api/model"
)

func TodoCreate(ctx *gin.Context) {
	// 将请求的 JSON 数据绑定到 TODO 结构体
	todo := model.TODO{}
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 获取 JWT token
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "JWT token is required",
		})
		return
	}

	// Validate JWT token
	claims, err := model.ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid JWT token",
		})
		return
	}
	todo.UserID = claims.UserID

	// 创建 TODO
	if err := model.CreateTODO(&todo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回创建完成的 TODO 信息
	ctx.JSON(http.StatusOK, gin.H{
		"todo": todo,
	})
}

func TodoGetService(ctx *gin.Context) {

}

func TodoUpdateService(ctx *gin.Context) {

}

func TodoDeleteService(ctx *gin.Context) {

}
