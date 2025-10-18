package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// asyncLoggerMiddleware demonstrates how to safely use gin.Context inside a goroutine.
func asyncLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Copy the context to use it inside the goroutine since the original
		// context is not goroutine-safe after the handler returns.
		copied := c.Copy()

		go func() {
			time.Sleep(2 * time.Second)
			log.Printf("[async] path=%s method=%s ip=%s", copied.Request.URL.Path, copied.Request.Method, copied.ClientIP())
		}()
	}
}

func main() {
	router := gin.Default()

	// Install the asynchronous logging middleware.
	router.Use(asyncLoggerMiddleware())

	// Simulate a long running handler with asynchronous logging.
	router.GET("/long_async", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"message": "long async task finished",
		})
	})

	// Simulate a synchronous handler to show that the middleware runs for every request.
	router.GET("/short_sync", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"message": "short sync task finished",
		})
	})

	log.Println("========================================")
	log.Println("Goroutines Inside Middleware 示例")
	log.Println("========================================")
	log.Println("服务器启动在 :8080")
	log.Println("")
	log.Println("测试命令：")
	log.Println("")
	log.Println("1. 触发异步日志：")
	log.Println("   curl http://localhost:8080/long_async")
	log.Println("")
	log.Println("2. 触发同步日志：")
	log.Println("   curl http://localhost:8080/short_sync")
	log.Println("")
	log.Println("观察日志输出，注意异步 goroutine 在请求返回后仍会运行。")
	log.Println("========================================")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
