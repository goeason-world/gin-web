package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的 Gin 引擎
	r := gin.Default()

	// 添加安全头中间件
	// 这个中间件会为所有响应添加安全相关的 HTTP 头
	r.Use(func(c *gin.Context) {
		// Content-Security-Policy (CSP)
		// 定义浏览器可以加载哪些资源，防止 XSS 攻击
		c.Header("Content-Security-Policy", "default-src 'self'")

		// X-Frame-Options
		// 防止网站被嵌入到 iframe 中，防止点击劫持攻击
		c.Header("X-Frame-Options", "DENY")

		// X-Content-Type-Options
		// 防止浏览器进行 MIME 类型嗅探，强制使用声明的 Content-Type
		c.Header("X-Content-Type-Options", "nosniff")

		// Strict-Transport-Security (HSTS)
		// 强制浏览器使用 HTTPS 连接，防止中间人攻击
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Referrer-Policy
		// 控制 Referer 头的发送策略，保护用户隐私
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy (原 Feature-Policy)
		// 控制浏览器功能和 API 的使用权限
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// 继续处理请求
		c.Next()
	})

	// 定义一个简单的路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 启动服务器
	r.Run(":8080")
}
