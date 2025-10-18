package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 定义一个 POST 路由 /form
	// 这个端点可以同时处理 urlencoded 和 multipart 两种编码方式
	router.POST("/form", func(c *gin.Context) {
		// c.PostForm() 用于获取表单字段的值
		// 如果字段不存在，返回空字符串
		message := c.PostForm("message")

		// c.DefaultPostForm() 用于获取表单字段的值
		// 如果字段不存在，返回指定的默认值
		nick := c.DefaultPostForm("nick", "anonymous")

		// 在控制台打印出接收到的数据，方便调试
		fmt.Printf("Message: %s\n", message)
		fmt.Printf("Nick: %s\n", nick)

		// 将接收到的数据以 JSON 格式返回给客户端
		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})

	// 打印测试说明
	println("=========================================")
	println("Multipart/Urlencoded 表单示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 演示如何处理不同编码方式的表单数据")
	println("- 支持 application/x-www-form-urlencoded")
	println("- 支持 multipart/form-data")
	println("- 使用 PostForm 和 DefaultPostForm 提取字段")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 URL 编码表单（application/x-www-form-urlencoded）：")
	println("   curl -X POST http://localhost:8080/form \\")
	println("     -H 'Content-Type: application/x-www-form-urlencoded' \\")
	println("     -d 'message=hello&nick=zhangsan'")
	println("")
	println("2. 测试 Multipart 表单（multipart/form-data）：")
	println("   curl -X POST http://localhost:8080/form \\")
	println("     -F 'message=world' \\")
	println("     -F 'nick=lisi'")
	println("")
	println("3. 测试默认值（不提供 nick 字段）：")
	println("   curl -X POST http://localhost:8080/form \\")
	println("     -H 'Content-Type: application/x-www-form-urlencoded' \\")
	println("     -d 'message=test'")
	println("")
	println("预期输出：")
	println("- 前两个请求会返回提供的值")
	println("- 第三个请求的 nick 字段会使用默认值 'anonymous'")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
