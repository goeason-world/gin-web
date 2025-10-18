package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginForm 定义了登录表单的数据结构
// 使用 form 标签来指定表单字段名
type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 定义一个 POST 路由 /login
	// 这个端点可以同时处理 urlencoded 和 multipart 两种编码方式
	router.POST("/login", func(c *gin.Context) {
		var form LoginForm

		// c.ShouldBind() 会根据 Content-Type 自动选择绑定方式：
		// - application/x-www-form-urlencoded: 绑定 URL 编码的表单数据
		// - multipart/form-data: 绑定 multipart 表单数据
		// 这使得我们可以用同一个处理函数处理两种编码方式
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 验证成功，返回用户信息
		c.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
			"user":    form.User,
		})
	})

	// 打印测试说明
	println("=========================================")
	println("Multipart/Urlencoded 绑定示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 演示如何处理不同编码方式的表单数据")
	println("- 支持 application/x-www-form-urlencoded")
	println("- 支持 multipart/form-data")
	println("- 使用同一个处理函数统一处理")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 URL 编码表单（application/x-www-form-urlencoded）：")
	println("   curl -X POST http://localhost:8080/login \\")
	println("     -H 'Content-Type: application/x-www-form-urlencoded' \\")
	println("     -d 'user=zhangsan&password=123456'")
	println("")
	println("2. 测试 Multipart 表单（multipart/form-data）：")
	println("   curl -X POST http://localhost:8080/login \\")
	println("     -F 'user=lisi' \\")
	println("     -F 'password=654321'")
	println("")
	println("3. 测试缺少必需字段（验证失败）：")
	println("   curl -X POST http://localhost:8080/login \\")
	println("     -H 'Content-Type: application/x-www-form-urlencoded' \\")
	println("     -d 'user=wangwu'")
	println("")
	println("预期输出：")
	println("- 前两个请求会返回成功响应")
	println("- 第三个请求会返回验证错误")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
