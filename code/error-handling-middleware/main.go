package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 是一个错误处理中间件，用于捕获并统一处理请求中的错误
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 步骤1: 先处理请求
		c.Next()

		// 步骤2: 检查上下文中是否有错误
		if len(c.Errors) > 0 {
			// 步骤3: 使用最后一个错误
			err := c.Errors.Last().Err

			// 步骤4: 返回统一的错误响应
			c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": err.Error(),
			})
		}

		// 如果没有错误，则不做任何处理
	}
}

func main() {
	r := gin.Default()

	// 挂载错误处理中间件
	r.Use(ErrorHandler())

	// 正常路由 - 不会触发错误
	r.GET("/ok", func(c *gin.Context) {
		somethingWentWrong := false

		if somethingWentWrong {
			c.Error(errors.New("something went wrong"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Everything is fine!",
		})
	})

	// 错误路由 - 会触发错误
	r.GET("/error", func(c *gin.Context) {
		somethingWentWrong := true

		if somethingWentWrong {
			c.Error(errors.New("something went wrong"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Everything is fine!",
		})
	})

	// 数据库错误示例
	r.GET("/db-error", func(c *gin.Context) {
		// 模拟数据库错误
		c.Error(errors.New("database connection failed"))
	})

	// 验证错误示例
	r.POST("/validate", func(c *gin.Context) {
		var data struct {
			Email string `json:"email"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.Error(errors.New("invalid request data"))
			return
		}

		// 简单的邮箱验证
		if data.Email == "" {
			c.Error(errors.New("email is required"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Validation passed",
		})
	})

	log.Println("========================================")
	log.Println("错误处理中间件示例")
	log.Println("========================================")
	log.Println("服务器启动在 :8080")
	log.Println("")
	log.Println("测试命令：")
	log.Println("1. 正常请求:")
	log.Println("   curl http://localhost:8080/ok")
	log.Println("")
	log.Println("2. 触发错误:")
	log.Println("   curl http://localhost:8080/error")
	log.Println("")
	log.Println("3. 数据库错误:")
	log.Println("   curl http://localhost:8080/db-error")
	log.Println("")
	log.Println("4. 验证错误:")
	log.Println("   curl -X POST http://localhost:8080/validate \\")
	log.Println("        -H \"Content-Type: application/json\" \\")
	log.Println("        -d '{}'")
	log.Println("========================================")
	log.Println("")

	r.Run(":8080")
}
