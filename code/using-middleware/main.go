package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// requestID 为每个请求生成一个简易的唯一 ID，并通过上下文和响应头暴露给后续逻辑。
func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := fmt.Sprintf("req-%d", time.Now().UnixNano())
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// requireAPIToken 用于演示路由组级别的中间件，校验 Authorization 头是否携带约定的令牌。
func requireAPIToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		const expected = "Bearer admin-token"
		if c.GetHeader("Authorization") != expected {
			// Abort 会阻断后续处理链，并立即返回响应。
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid token",
			})
			return
		}

		c.Next()
	}
}

// measureLatency 展示按路由挂载中间件：记录开始时间并在响应头中加入 X-Response-Time。
func measureLatency() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 下游 handler 可以读取 request_start 实现更细粒度的统计。
		c.Set("request_start", start)

		c.Next()

		elapsed := time.Since(start)
		c.Header("X-Response-Time", fmt.Sprintf("%dms", elapsed.Milliseconds()))
	}
}

func main() {
	// 使用 gin.New 自行挂载全局中间件，避免重复添加默认 Logger。
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), requestID())

	// 全局中间件已经设置 request_id，这里直接读取并返回给客户端。
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":    "pong",
			"request_id": c.GetString("request_id"),
		})
	})

	// 对 /admin 路由组新增鉴权中间件。
	admin := router.Group("/admin")
	admin.Use(requireAPIToken())
	admin.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":     "ok",
			"request_id": c.GetString("request_id"),
		})
	})

	// 仅在 /reports/daily 上挂载 measureLatency，展示按路由使用中间件的写法。
	router.GET("/reports/daily", measureLatency(), func(c *gin.Context) {
		start := c.MustGet("request_start").(time.Time)
		// 模拟耗时计算。
		time.Sleep(120 * time.Millisecond)

		c.JSON(http.StatusOK, gin.H{
			"generated_at": start.Format(time.RFC3339),
			"request_id":   c.GetString("request_id"),
		})
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
