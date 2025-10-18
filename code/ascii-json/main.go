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
		// 准备要返回的数据，包含中文和 HTML 标签
		data := map[string]interface{}{
			"lang": "GO语言", // 中文字符
			"tag":  "<br>", // HTML 标签
		}

		// 使用 AsciiJSON 返回响应
		// 所有非 ASCII 字符都会被转义为 \uXXXX 格式
		// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	// 打印测试说明
	println("========================================")
	println("AsciiJSON 示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 AsciiJSON：")
	println("   curl http://localhost:8080/someJSON")
	println("")
	println("预期输出：")
	println(`   {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}`)
	println("")
	println("说明：")
	println("- 中文字符被转义为 Unicode 转义序列")
	println("- HTML 标签被转义，防止 XSS 攻击")
	println("========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
