package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	// 使用 JSON 格式的日志
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 构建日志对象
		log := map[string]interface{}{
			"timestamp":   param.TimeStamp.Format("2006-01-02T15:04:05Z07:00"),
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"protocol":    param.Request.Proto,
			"status_code": param.StatusCode,
			"latency":     param.Latency.String(),
			"latency_ms":  param.Latency.Milliseconds(),
			"user_agent":  param.Request.UserAgent(),
			"error":       param.ErrorMessage,
		}

		// 转换为 JSON 字符串
		jsonLog, _ := json.Marshal(log)
		return string(jsonLog) + "\n"
	}))

	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, JSON log format!",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello, %s!", name),
		})
	})

	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
	})

	fmt.Println("========================================")
	fmt.Println("JSON 格式日志示例")
	fmt.Println("========================================")
	fmt.Println("服务器启动在 :8080")
	fmt.Println("")
	fmt.Println("测试命令：")
	fmt.Println("1. curl http://localhost:8080/")
	fmt.Println("2. curl http://localhost:8080/ping")
	fmt.Println("3. curl http://localhost:8080/user/张三")
	fmt.Println("4. curl http://localhost:8080/error")
	fmt.Println("========================================")
	fmt.Println("")

	router.Run(":8080")
}
