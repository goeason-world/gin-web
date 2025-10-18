package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 默认情况下，Gin 会根据输出终端自动决定是否使用彩色日志
	// 如果输出到终端（TTY），则使用彩色日志
	// 如果输出到文件或管道，则使用纯文本日志

	// 禁用日志颜色
	// 当你需要将日志输出到文件或日志收集系统时，应该禁用颜色
	// gin.DisableConsoleColor()

	// 强制启用日志颜色
	// 即使输出不是终端，也会使用 ANSI 颜色代码
	gin.ForceConsoleColor()

	router := gin.Default()

	// 定义一些测试路由
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello " + name,
		})
	})

	router.POST("/user", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created",
		})
	})

	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
	})

	// 打印测试说明
	println("========================================")
	println("控制日志输出颜色示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("当前设置：强制启用彩色日志")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 200 状态码（绿色）：")
	println(`   curl http://localhost:8080/ping`)
	println("")
	println("2. 测试带参数的路由（绿色）：")
	println(`   curl http://localhost:8080/hello/张三`)
	println("")
	println("3. 测试 POST 请求（绿色）：")
	println(`   curl -X POST http://localhost:8080/user`)
	println("")
	println("4. 测试 500 错误（红色）：")
	println(`   curl http://localhost:8080/error`)
	println("")
	println("5. 测试 404 错误（黄色）：")
	println(`   curl http://localhost:8080/notfound`)
	println("")
	println("提示：")
	println("- 在代码中取消注释 gin.DisableConsoleColor() 可以禁用颜色")
	println("- 在代码中注释掉 gin.ForceConsoleColor() 可以使用自动检测")
	println("========================================")
	println("")

	// 监听 8089 端口
	router.Run(":8080")
}
