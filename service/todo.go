package service

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"todo-list-api/middleware"
	"todo-list-api/model"
)

func TodoCreate(ctx *gin.Context) {
	// 将请求的 JSON 数据绑定到 TODO 结构体
	todo := model.TODO{}
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	// 以下是TodoCreate函数中的修改部分，将model.ValidateToken替换为middleware.ValidateJWT

// 获取 JWT token
token := ctx.GetHeader("Authorization")
if token == "" {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "message": "Unauthorized",
    })
    return
}

// 从Authorization头获取JWT部分
parts := strings.Split(token, " ")
if len(parts) == 2 && parts[0] == "Bearer" {
    token = parts[1]
} else {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "message": "Invalid token format",
    })
    return
}

// 验证JWT token
claims, err := middleware.ValidateJWT(token)
if err != nil {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "message": "Unauthorized",
    })
    return
}
todo.UserID = claims.UserID
	// 创建 TODO
	if err := model.CreateTODO(&todo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// 返回创建完成的 TODO 信息
	ctx.JSON(http.StatusCreated, gin.H{
		"id":          todo.ID,
		"title":       todo.Title,
		"description": todo.Description,
	})
}

// TodoGet 获取待办事项列表
func TodoGet(ctx *gin.Context) {
	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// 获取筛选参数
	filter := make(map[string]interface{})
	if completed, exists := ctx.GetQuery("completed"); exists {
		if completed == "true" {
			filter["completed"] = true
		} else if completed == "false" {
			filter["completed"] = false
		}
	}

	if title := ctx.Query("title"); title != "" {
		filter["title"] = title
	}

	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		// 获取 JWT token
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		// Validate JWT token
		claims, err := model.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		userID = claims.UserID
	}

	// 查询待办事项
	todos, count, err := model.GetTODOsByUserID(userID.(int), page, limit, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// 计算总页数
	totalPages := (count + int64(limit) - 1) / int64(limit)

	ctx.JSON(http.StatusOK, gin.H{
		"todos": todos,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"limit":        limit,
			"total_items":  count,
		},
	})
}

// TodoGetByID 获取单个待办事项
func TodoGetByID(ctx *gin.Context) {
	// 从路径获取待办事项ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid todo ID",
		})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		// 获取 JWT token
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		// Validate JWT token
		claims, err := model.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		userID = claims.UserID
	}

	// 获取待办事项
	todo, err := model.GetTODOByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Todo not found",
		})
		return
	}

	// 检查用户是否有权限查看
	if todo.UserID != userID.(int) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          todo.ID,
		"title":       todo.Title,
		"description": todo.Description,
		"completed":   todo.Completed,
	})
}

// TodoUpdate 更新待办事项
func TodoUpdate(ctx *gin.Context) {
	// 从路径获取待办事项ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid todo ID",
		})
		return
	}

	// 验证请求数据
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		// 获取 JWT token
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		// Validate JWT token
		claims, err := model.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		userID = claims.UserID
	}

	// 检查用户是否有权限更新
	if !model.IsOwner(id, userID.(int)) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden",
		})
		return
	}

	// 获取并更新待办事项
	todo, err := model.GetTODOByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Todo not found",
		})
		return
	}

	todo.Title = req.Title
	todo.Description = req.Description

	if err := model.UpdateTODO(&todo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// 返回更新后的待办事项
	ctx.JSON(http.StatusOK, gin.H{
		"id":          todo.ID,
		"title":       todo.Title,
		"description": todo.Description,
	})
}

// TodoDelete 删除待办事项
func TodoDelete(ctx *gin.Context) {
	// 从路径获取待办事项ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid todo ID",
		})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		// 获取 JWT token
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		// Validate JWT token
		claims, err := model.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		userID = claims.UserID
	}

	// 检查用户是否有权限删除
	if !model.IsOwner(id, userID.(int)) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden",
		})
		return
	}

	// 删除待办事项
	if err := model.DeleteTODOByID(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// TodoComplete 标记待办事项为完成/未完成
func TodoComplete(ctx *gin.Context) {
	// 从路径获取待办事项ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid todo ID",
		})
		return
	}

	// 获取完成状态
	type CompleteRequest struct {
		Completed bool `json:"completed"`
	}
	var req CompleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		// 获取 JWT token
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		// Validate JWT token
		claims, err := model.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		userID = claims.UserID
	}

	// 检查用户是否有权限更新
	if !model.IsOwner(id, userID.(int)) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden",
		})
		return
	}

	// 获取并更新待办事项
	todo, err := model.GetTODOByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Todo not found",
		})
		return
	}

	todo.Completed = req.Completed

	if err := model.UpdateTODO(&todo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          todo.ID,
		"title":       todo.Title,
		"description": todo.Description,
		"completed":   todo.Completed,
	})
}
