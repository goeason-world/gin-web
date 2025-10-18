package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 使用 c.JSON() - 默认方法，会转义 HTML 特殊字符
	router.GET("/json", func(c *gin.Context) {
		// 包含 HTML 特殊字符的数据
		c.JSON(http.StatusOK, gin.H{
			"html": "<b>Hello, World!</b>",
			"url":  "https://example.com?a=1&b=2",
		})
		// 输出: {"html":"\u003cb\u003eHello, World!\u003c/b\u003e","url":"https://example.com?a=1\u0026b=2"}
		// 注意: < > & 被转义为 \u003c \u003e \u0026
	})

	// 使用 c.PureJSON() - 不转义 HTML 特殊字符
	router.GET("/purejson", func(c *gin.Context) {
		// 包含 HTML 特殊字符的数据
		c.PureJSON(http.StatusOK, gin.H{
			"html": "<b>Hello, World!</b>",
			"url":  "https://example.com?a=1&b=2",
		})
		// 输出: {"html":"<b>Hello, World!</b>","url":"https://example.com?a=1&b=2"}
		// 注意: < > & 保持原样，不被转义
	})

	// 使用 c.AsciiJSON() - 转义所有非 ASCII 字符
	router.GET("/asciijson", func(c *gin.Context) {
		// 包含中文和 HTML 特殊字符的数据
		c.AsciiJSON(http.StatusOK, gin.H{
			"html":    "<b>Hello, World!</b>",
			"chinese": "你好，世界！",
		})
		// 输出: {"chinese":"\u4f60\u597d\uff0c\u4e16\u754c\uff01","html":"\u003cb\u003eHello, World!\u003c/b\u003e"}
		// 注意: 中文字符和 HTML 特殊字符都被转义
	})

	// 对比示例：展示三种方法的区别
	router.GET("/compare", func(c *gin.Context) {
		data := gin.H{
			"message": "比较三种 JSON 方法",
			"html":    "<script>alert('XSS')</script>",
			"chinese": "中文测试",
			"url":     "https://example.com?foo=bar&baz=qux",
		}

		// 根据查询参数选择不同的方法
		method := c.DefaultQuery("method", "json")

		switch method {
		case "pure":
			c.PureJSON(http.StatusOK, data)
		case "ascii":
			c.AsciiJSON(http.StatusOK, data)
		default:
			c.JSON(http.StatusOK, data)
		}
	})

	// 打印测试说明
	println("=========================================")
	println("PureJSON 示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- c.JSON(): 转义 HTML 特殊字符（默认）")
	println("- c.PureJSON(): 不转义 HTML 特殊字符")
	println("- c.AsciiJSON(): 转义所有非 ASCII 字符")
	println("")
	println("测试命令：")
	println("")
	println("1. 使用 c.JSON()（默认方法）：")
	println("   curl http://localhost:8080/json")
	println("")
	println("2. 使用 c.PureJSON()（不转义 HTML）：")
	println("   curl http://localhost:8080/purejson")
	println("")
	println("3. 使用 c.AsciiJSON()（转义所有非 ASCII）：")
	println("   curl http://localhost:8080/asciijson")
	println("")
	println("4. 对比三种方法：")
	println("   curl 'http://localhost:8080/compare?method=json'")
	println("   curl 'http://localhost:8080/compare?method=pure'")
	println("   curl 'http://localhost:8080/compare?method=ascii'")
	println("")
	println("预期输出：")
	println("- JSON: HTML 字符被转义为 \\uXXXX")
	println("- PureJSON: HTML 字符保持原样")
	println("- AsciiJSON: 所有非 ASCII 字符都被转义")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
