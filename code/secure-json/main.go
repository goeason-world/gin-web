package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 定义一个 GET 路由
	router.GET("/someJSON", func(c *gin.Context) {
		// 准备要返回的数据
		names := []string{"lena", "austin", "foo"}

		// 使用 SecureJSON 返回响应
		// SecureJSON 会在 JSON 前添加 "while(1);" 前缀，防止 JSON 劫持攻击
		// 输出 : while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})

	// 定义一个自定义前缀的路由
	router.GET("/customPrefix", func(c *gin.Context) {
		// 准备要返回的数据
		data := gin.H{
			"user":  "alice",
			"email": "alice@example.com",
		}

		// 使用自定义前缀的 SecureJSON
		// 你也可以使用自己的 SecureJSON 前缀
		// 输出 : )]}',\n{"email":"alice@example.com","user":"alice"}
		c.SecureJSON(http.StatusOK, data)
	})

	// 打印测试说明
	println("========================================")
	println("SecureJSON 示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试默认 SecureJSON：")
	println("   curl http://localhost:8080/someJSON")
	println("")
	println("预期输出：")
	println(`   while(1);["lena","austin","foo"]`)
	println("")
	println("2. 测试自定义前缀的 SecureJSON：")
	println("   curl http://localhost:8080/customPrefix")
	println("")
	println("预期输出：")
	println(`   while(1);{"email":"alice@example.com","user":"alice"}`)
	println("")
	println("说明：")
	println("- SecureJSON 在 JSON 前添加 'while(1);' 前缀")
	println("- 这可以防止 JSON 劫持攻击（JSON Hijacking）")
	println("- 客户端需要在解析前移除这个前缀")
	println("========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
