package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// loginEndpoint 是一个模拟的登录处理函数
func loginEndpoint(c *gin.Context) {
	// 在实际应用中，这里会处理登录逻辑
	// 比如验证用户名和密码
	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"path":    c.Request.URL.Path,
	})
}

// submitEndpoint 是一个模拟的提交数据处理函数
func submitEndpoint(c *gin.Context) {
	// 在实际应用中，这里会处理数据提交
	c.JSON(http.StatusOK, gin.H{
		"message": "submit successful",
		"path":    c.Request.URL.Path,
	})
}

// readEndpoint 是一个模拟的读取数据处理函数
func readEndpoint(c *gin.Context) {
	// 在实际应用中，这里会处理数据读取
	c.JSON(http.StatusOK, gin.H{
		"message": "read successful",
		"path":    c.Request.URL.Path,
	})
}

func main() {
	// 创建一个默认的 Gin 引擎
	router := gin.Default()

	// === 路由组 v1 ===
	// 使用 router.Group() 创建一个名为 "/v1" 的路由组
	// 所有在这个组内定义的路由，其 URL 都会自动加上 "/v1" 前缀
	v1 := router.Group("/v1")
	{
		// 注册 POST /v1/login 路由
		// 实际访问路径是 /v1/login
		v1.POST("/login", loginEndpoint)

		// 注册 POST /v1/submit 路由
		// 实际访问路径是 /v1/submit
		v1.POST("/submit", submitEndpoint)

		// 注册 POST /v1/read 路由
		// 实际访问路径是 /v1/read
		v1.POST("/read", readEndpoint)
	}

	// === 路由组 v2 ===
	// 同样地，创建一个名为 "/v2" 的路由组
	v2 := router.Group("/v2")
	{
		// 注册 POST /v2/login 路由
		// 实际访问路径是 /v2/login
		v2.POST("/login", loginEndpoint)

		// 注册 POST /v2/submit 路由
		// 实际访问路径是 /v2/submit
		v2.POST("/submit", submitEndpoint)

		// 注册 POST /v2/read 路由
		// 实际访问路径是 /v2/read
		v2.POST("/read", readEndpoint)
	}
	
	// 打印测试说明
	println("==========================================")
	println("路由分组示例")
	println("==========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令:")
	println("")
	println("1. 测试 v1 组的 /login:")
	println("   curl -X POST http://localhost:8080/v1/login")
	println("")
	println("2. 测试 v1 组的 /submit:")
	println("   curl -X POST http://localhost:8080/v1/submit")
	println("")
	println("3. 测试 v2 组的 /login:")
	println("   curl -X POST http://localhost:8080/v2/login")
	println("")
	println("==========================================")
	println("")

	// 启动服务器，监听在 8080 端口
	router.Run(":8080")
}