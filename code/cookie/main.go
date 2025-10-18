package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 路由器
	router := gin.Default()

	// 路由 1: 获取和设置 Cookie
	// 演示如何读取 Cookie，如果不存在则创建一个新的 Cookie
	router.GET("/cookie", func(c *gin.Context) {
		// 尝试获取名为 "gin_cookie" 的 Cookie
		cookie, err := c.Cookie("gin_cookie")

		// 如果 Cookie 不存在（首次访问）
		if err != nil {
			cookie = "NotSet"
			// 设置 Cookie
			// 参数说明：
			// - name: Cookie 名称 ("gin_cookie")
			// - value: Cookie 值 ("test")
			// - maxAge: Cookie 有效期（秒），3600 表示 1 小时
			// - path: Cookie 的路径作用域，"/" 表示整个域名下都有效
			// - domain: Cookie 的域名作用域
			// - secure: 是否只在 HTTPS 下传输，false 表示 HTTP 和 HTTPS 都可以
			// - httpOnly: 是否禁止 JavaScript 访问，true 表示只能通过 HTTP(S) 访问，增强安全性
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		// 在服务器控制台打印 Cookie 值
		fmt.Printf("Cookie value: %s \n", cookie)

		// 返回 JSON 响应
		c.JSON(http.StatusOK, gin.H{
			"message": "Cookie 操作成功",
			"cookie":  cookie,
		})
	})

	// 路由 2: 更新 Cookie
	// 演示如何更新现有的 Cookie 值
	router.GET("/update-cookie", func(c *gin.Context) {
		// 设置新的 Cookie 值
		c.SetCookie("gin_cookie", "updated_value", 3600, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{
			"message": "Cookie 已更新",
			"value":   "updated_value",
		})
	})

	// 路由 3: 删除 Cookie
	// 演示如何删除 Cookie（通过设置 maxAge 为 -1）
	router.GET("/delete-cookie", func(c *gin.Context) {
		// 删除 Cookie 的关键是将 maxAge 设置为 -1
		// 这会告诉浏览器立即删除该 Cookie
		c.SetCookie("gin_cookie", "", -1, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{
			"message": "Cookie 已删除",
		})
	})

	// 路由 4: 设置多个 Cookie
	// 演示如何同时设置多个不同的 Cookie
	router.GET("/multiple-cookies", func(c *gin.Context) {
		// 设置用户信息 Cookie
		c.SetCookie("user_id", "12345", 3600, "/", "localhost", false, true)
		// 设置会话 Cookie
		c.SetCookie("session_token", "abc123xyz", 7200, "/", "localhost", false, true)
		// 设置偏好设置 Cookie（允许 JavaScript 访问）
		c.SetCookie("theme", "dark", 86400, "/", "localhost", false, false)

		c.JSON(http.StatusOK, gin.H{
			"message": "已设置多个 Cookie",
			"cookies": []string{"user_id", "session_token", "theme"},
		})
	})

	// 路由 5: 读取所有 Cookie
	// 演示如何读取请求中的所有 Cookie
	router.GET("/all-cookies", func(c *gin.Context) {
		// 获取所有 Cookie
		cookies := c.Request.Cookies()

		// 构建 Cookie 信息映射
		cookieMap := make(map[string]string)
		for _, cookie := range cookies {
			cookieMap[cookie.Name] = cookie.Value
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "所有 Cookie",
			"count":   len(cookies),
			"cookies": cookieMap,
		})
	})

	// 打印测试说明
	fmt.Println("========================================")
	fmt.Println("Cookie 操作示例")
	fmt.Println("========================================")
	fmt.Println("服务器启动在 :8080")
	fmt.Println("")
	fmt.Println("测试命令：")
	fmt.Println("")
	fmt.Println("1. 获取/设置 Cookie（首次访问）：")
	fmt.Println("   curl -c cookies.txt http://localhost:8080/cookie")
	fmt.Println("")
	fmt.Println("2. 使用已保存的 Cookie 再次访问：")
	fmt.Println("   curl -b cookies.txt http://localhost:8080/cookie")
	fmt.Println("")
	fmt.Println("3. 更新 Cookie：")
	fmt.Println("   curl -b cookies.txt -c cookies.txt http://localhost:8080/update-cookie")
	fmt.Println("")
	fmt.Println("4. 删除 Cookie：")
	fmt.Println("   curl -b cookies.txt -c cookies.txt http://localhost:8080/delete-cookie")
	fmt.Println("")
	fmt.Println("5. 设置多个 Cookie：")
	fmt.Println("   curl -c cookies.txt http://localhost:8080/multiple-cookies")
	fmt.Println("")
	fmt.Println("6. 读取所有 Cookie：")
	fmt.Println("   curl -b cookies.txt http://localhost:8080/all-cookies")
	fmt.Println("")
	fmt.Println("========================================")
	fmt.Println("")

	// 启动服务器，监听 8080 端口
	router.Run(":8080")
}
