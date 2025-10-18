package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 自定义错误类型
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

// 增强的错误处理中间件
func EnhancedErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 根据错误类型设置不同的状态码
			statusCode := http.StatusInternalServerError
			message := err.Error()

			// 检查是否是自定义错误
			var appErr *AppError
			if errors.As(err, &appErr) {
				statusCode = appErr.Code
				message = appErr.Message
			}

			c.JSON(statusCode, map[string]any{
				"success": false,
				"message": message,
			})
		}
	}
}

func main() {
	r := gin.Default()

	r.Use(EnhancedErrorHandler())

	r.GET("/not-found", func(c *gin.Context) {
		c.Error(&AppError{
			Code:    http.StatusNotFound,
			Message: "Resource not found",
		})
	})

	r.GET("/unauthorized", func(c *gin.Context) {
		c.Error(&AppError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized access",
		})
	})

	r.GET("/bad-request", func(c *gin.Context) {
		c.Error(&AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request parameters",
		})
	})

	r.GET("/server-error", func(c *gin.Context) {
		// 普通错误，会使用默认的 500 状态码
		c.Error(errors.New("internal server error"))
	})

	log.Println("========================================")
	log.Println("增强的错误处理中间件示例")
	log.Println("========================================")
	log.Println("服务器启动在 :8080")
	log.Println("")
	log.Println("测试命令：")
	log.Println("1. 404 错误:")
	log.Println("   curl http://localhost:8080/not-found")
	log.Println("")
	log.Println("2. 401 未授权:")
	log.Println("   curl http://localhost:8080/unauthorized")
	log.Println("")
	log.Println("3. 400 错误请求:")
	log.Println("   curl http://localhost:8080/bad-request")
	log.Println("")
	log.Println("4. 500 服务器错误:")
	log.Println("   curl http://localhost:8080/server-error")
	log.Println("========================================")
	log.Println("")

	r.Run(":8080")
}
