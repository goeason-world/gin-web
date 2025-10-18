package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 1. HTTP 重定向 - 外部 URL
	// 使用 302 临时重定向到外部网站
	router.GET("/external", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://www.google.com")
	})

	// 2. HTTP 重定向 - 内部路由
	// 使用 301 永久重定向到内部路由
	router.GET("/old-path", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/new-path")
	})

	router.GET("/new-path", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "这是新路径",
		})
	})

	// 3. POST 后重定向到 GET（PRG 模式）
	// 使用 303 状态码，防止表单重复提交
	router.POST("/submit", func(c *gin.Context) {
		// 处理表单提交
		name := c.PostForm("name")

		fmt.Printf("Received submission: %s\n", name)

		// 保存数据到数据库...
		// db.Save(name)

		// 重定向到成功页面
		c.Redirect(http.StatusSeeOther, "/success")
	})

	router.GET("/success", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "提交成功",
		})
	})

	// 4. 路由重定向（内部重定向）
	// 不会发送 HTTP 重定向响应，而是在服务器内部处理
	router.GET("/test", func(c *gin.Context) {
		// 修改请求路径
		c.Request.URL.Path = "/test2"
		// 让 Gin 重新处理请求
		router.HandleContext(c)
	})

	router.GET("/test2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "通过路由重定向到达这里",
		})
	})

	// 5. 条件重定向
	// 根据条件决定是否重定向
	router.GET("/conditional", func(c *gin.Context) {
		// 检查用户是否登录
		isLoggedIn := c.Query("logged_in") == "true"

		if !isLoggedIn {
			// 未登录，重定向到登录页
			c.Redirect(http.StatusFound, "/login")
			return
		}

		// 已登录，显示内容
		c.JSON(http.StatusOK, gin.H{
			"message": "欢迎回来",
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "请登录",
		})
	})

	// 6. 带查询参数的重定向
	router.GET("/redirect-with-params", func(c *gin.Context) {
		// 重定向时保留或添加查询参数
		c.Redirect(http.StatusFound, "/target?source=redirect&timestamp=123456")
	})

	router.GET("/target", func(c *gin.Context) {
		source := c.Query("source")
		timestamp := c.Query("timestamp")

		c.JSON(http.StatusOK, gin.H{
			"message":   "到达目标页面",
			"source":    source,
			"timestamp": timestamp,
		})
	})

	// 7. 短链接服务示例
	// 模拟短链接重定向
	shortLinks := map[string]string{
		"google": "https://www.google.com",
		"github": "https://github.com",
		"gin":    "https://gin-gonic.com",
	}

	router.GET("/s/:code", func(c *gin.Context) {
		code := c.Param("code")

		// 查找短链接对应的原始 URL
		if url, exists := shortLinks[code]; exists {
			// 使用 301 永久重定向
			c.Redirect(http.StatusMovedPermanently, url)
			return
		}

		// 短链接不存在
		c.JSON(http.StatusNotFound, gin.H{
			"error": "短链接不存在",
		})
	})

	// 8. 语言重定向
	// 根据用户语言偏好重定向
	router.GET("/", func(c *gin.Context) {
		// 获取用户首选语言
		acceptLanguage := c.GetHeader("Accept-Language")

		// 简单的语言检测
		if len(acceptLanguage) >= 2 {
			lang := acceptLanguage[:2]
			switch lang {
			case "zh":
				c.Redirect(http.StatusFound, "/zh/home")
				return
			case "en":
				c.Redirect(http.StatusFound, "/en/home")
				return
			}
		}

		// 默认重定向到英文页面
		c.Redirect(http.StatusFound, "/en/home")
	})

	router.GET("/zh/home", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "欢迎来到中文主页",
		})
	})

	router.GET("/en/home", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to English homepage",
		})
	})

	// 9. HTTPS 重定向
	// 强制使用 HTTPS
	router.GET("/secure", func(c *gin.Context) {
		// 检查是否使用 HTTPS
		if c.Request.TLS == nil && c.GetHeader("X-Forwarded-Proto") != "https" {
			// 构建 HTTPS URL
			httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusMovedPermanently, httpsURL)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "这是安全连接",
		})
	})

	// 10. 移动端重定向
	// 根据 User-Agent 重定向到移动版
	router.GET("/page", func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")

		// 简单的移动设备检测
		isMobile := false
		mobileKeywords := []string{"Mobile", "Android", "iPhone", "iPad"}
		for _, keyword := range mobileKeywords {
			if contains(userAgent, keyword) {
				isMobile = true
				break
			}
		}

		if isMobile {
			c.Redirect(http.StatusFound, "/m/page")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "桌面版页面",
		})
	})

	router.GET("/m/page", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "移动版页面",
		})
	})

	// 打印测试说明
	println("=========================================")
	println("重定向示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- HTTP 重定向：使用 c.Redirect()")
	println("- 路由重定向：使用 HandleContext()")
	println("- 支持多种重定向状态码")
	println("- 条件重定向和参数传递")
	println("")
	println("测试命令：")
	println("")
	println("1. 外部重定向：")
	println("   curl -L http://localhost:8080/external")
	println("")
	println("2. 内部重定向：")
	println("   curl -L http://localhost:8080/old-path")
	println("")
	println("3. POST 后重定向：")
	println("   curl -L -X POST http://localhost:8080/submit -d 'name=test'")
	println("")
	println("4. 路由重定向：")
	println("   curl http://localhost:8080/test")
	println("")
	println("5. 条件重定向（未登录）：")
	println("   curl -L http://localhost:8080/conditional")
	println("")
	println("6. 条件重定向（已登录）：")
	println("   curl http://localhost:8080/conditional?logged_in=true")
	println("")
	println("7. 短链接重定向：")
	println("   curl -L http://localhost:8080/s/google")
	println("")
	println("8. 语言重定向：")
	println("   curl -L -H 'Accept-Language: zh-CN' http://localhost:8080/")
	println("")
	println("注意：")
	println("- 使用 -L 参数让 curl 自动跟随重定向")
	println("- 不使用 -L 可以看到重定向响应头")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}

// 辅助函数：检查字符串是否包含子串
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
