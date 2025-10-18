package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的 Gin 引擎
	r := gin.Default()

	// 定义一个 GET 路由 /JSONP
	// 这个端点将演示 JSONP 的工作方式
	r.GET("/JSONP", func(c *gin.Context) {
		// 准备要返回的 JSON 数据
		data := gin.H{
			"foo": "bar",
		}

		// 使用 c.JSONP() 返回响应
		// Gin 会自动检查 URL 中是否有名为 "callback" 的查询参数
		// 如果有，它会将 data 包裹在回调函数中返回
		// 如果没有，它会返回一个普通的 JSON 响应
		c.JSONP(http.StatusOK, data)
	})

	// 打印测试说明
	println("=========================================")
	println("JSONP 示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试普通的 JSON 响应（不带 callback）：")
	println("   curl http://localhost:8080/JSONP")
	println("")
	println("2. 测试 JSONP 响应（带 callback）：")
	println("   curl http://localhost:8080/JSONP?callback=myCallbackFunction")
	println("")
	println("预期输出：")
	println("1. 普通 JSON: {\"foo\":\"bar\"}")
	println("2. JSONP:     myCallbackFunction({\"foo\":\"bar\"});")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
