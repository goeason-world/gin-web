package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 是一个自定义的日志中间件
// 它会记录每个请求的详细信息
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		startTime := time.Now()

		// 记录请求信息
		log.Printf("开始处理请求: %s %s", c.Request.Method, c.Request.URL.Path)

		// 处理请求
		c.Next()

		// 请求结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 记录响应信息
		log.Printf("请求处理完成: %s %s | 状态码: %d | 耗时: %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)
	}
}

// AuthRequired 是一个简单的身份验证中间件
// 它检查请求头中是否包含有效的 token
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证令牌",
			})
			c.Abort() // 终止请求，不再执行后续的处理函数
			return
		}

		// 简单的 token 验证（实际应用中应该验证 JWT 或查询数据库）
		if token != "Bearer valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的认证令牌",
			})
			c.Abort()
			return
		}

		// 验证通过，可以在上下文中设置用户信息
		c.Set("user_id", "12345")
		c.Set("username", "张三")

		// 继续处理请求
		c.Next()
	}
}

// CORS 是一个跨域资源共享中间件
// 它添加必要的 CORS 响应头
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimiter 是一个简单的限流中间件
// 它限制每个 IP 的请求频率
func RateLimiter() gin.HandlerFunc {
	// 存储每个 IP 的最后请求时间
	lastRequestTime := make(map[string]time.Time)
	minInterval := 1 * time.Second // 最小请求间隔

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		if lastTime, exists := lastRequestTime[clientIP]; exists {
			if now.Sub(lastTime) < minInterval {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "请求过于频繁，请稍后再试",
				})
				c.Abort()
				return
			}
		}

		lastRequestTime[clientIP] = now
		c.Next()
	}
}

func main() {
	// 创建不带默认中间件的 Gin 引擎
	router := gin.New()

	// 使用 Recovery 中间件（捕获 panic）
	router.Use(gin.Recovery())

	// 使用自定义的日志中间件
	router.Use(Logger())

	// 使用 CORS 中间件
	router.Use(CORS())

	// 公开路由（不需要认证）
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "欢迎访问公开 API",
		})
	})

	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "这是一个公开的端点",
		})
	})

	// 需要认证的路由组
	authorized := router.Group("/api")
	authorized.Use(AuthRequired()) // 应用认证中间件
	{
		authorized.GET("/profile", func(c *gin.Context) {
			// 从上下文中获取用户信息
			userID, _ := c.Get("user_id")
			username, _ := c.Get("username")

			c.JSON(http.StatusOK, gin.H{
				"message":  "用户资料",
				"user_id":  userID,
				"username": username,
			})
		})

		authorized.GET("/data", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "这是受保护的数据",
				"data":    []string{"item1", "item2", "item3"},
			})
		})
	}

	// 应用限流中间件的路由
	limited := router.Group("/limited")
	limited.Use(RateLimiter())
	{
		limited.GET("/resource", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "这是一个限流的资源",
			})
		})
	}

	// 打印测试说明
	fmt.Println("========================================")
	fmt.Println("自定义中间件示例")
	fmt.Println("========================================")
	fmt.Println("服务器启动在 :8080")
	fmt.Println("")
	fmt.Println("测试命令：")
	fmt.Println("")
	fmt.Println("1. 测试公开路由：")
	fmt.Println("   curl http://localhost:8080/")
	fmt.Println("   curl http://localhost:8080/public")
	fmt.Println("")
	fmt.Println("2. 测试需要认证的路由（无 token）：")
	fmt.Println("   curl http://localhost:8080/api/profile")
	fmt.Println("")
	fmt.Println("3. 测试需要认证的路由（有效 token）：")
	fmt.Println("   curl -H \"Authorization: Bearer valid-token\" http://localhost:8080/api/profile")
	fmt.Println("   curl -H \"Authorization: Bearer valid-token\" http://localhost:8080/api/data")
	fmt.Println("")
	fmt.Println("4. 测试限流路由（快速连续请求）：")
	fmt.Println("   curl http://localhost:8080/limited/resource")
	fmt.Println("   curl http://localhost:8080/limited/resource  # 立即再次请求")
	fmt.Println("")
	fmt.Println("========================================")
	fmt.Println("")

	// 启动服务器
	router.Run(":8080")
}
