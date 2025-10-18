package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个不带任何中间件的 Gin 引擎
	router := gin.New()

	// 使用自定义的日志格式
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
		// 你可以根据需要格式化日志输出
		return fmt.Sprintf("[%s] - %s \"%s %s %s %d %s \"%s\" %s\"\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"), // 时间戳
			param.ClientIP,            // 客户端 IP
			param.Method,              // HTTP 方法
			param.Path,                // 请求路径
			param.Request.Proto,       // HTTP 协议版本
			param.StatusCode,          // 状态码
			param.Latency,             // 延迟时间
			param.Request.UserAgent(), // User-Agent
			param.ErrorMessage,        // 错误信息
		)
	}))

	// 使用 Recovery 中间件
	router.Use(gin.Recovery())

	// 定义路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, custom log format!",
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

	// 打印测试说明
	fmt.Println("========================================")
	fmt.Println("自定义日志格式示例")
	fmt.Println("========================================")
	fmt.Println("服务器启动在 :8080")
	fmt.Println("")
	fmt.Println("日志格式说明：")
	fmt.Println("[时间戳] - 客户端IP \"方法 路径 协议 状态码 延迟 User-Agent 错误信息\"")
	fmt.Println("")
	fmt.Println("测试命令：")
	fmt.Println("")
	fmt.Println("1. 测试根路径：")
	fmt.Println("   curl http://localhost:8080/")
	fmt.Println("")
	fmt.Println("2. 测试 ping：")
	fmt.Println("   curl http://localhost:8080/ping")
	fmt.Println("")
	fmt.Println("3. 测试带参数的路由：")
	fmt.Println("   curl http://localhost:8080/user/张三")
	fmt.Println("")
	fmt.Println("4. 测试错误响应：")
	fmt.Println("   curl http://localhost:8080/error")
	fmt.Println("")
	fmt.Println("========================================")
	fmt.Println("")

	// 启动服务器
	router.Run(":8080")
}
