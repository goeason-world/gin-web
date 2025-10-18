package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// User 结构体用于演示各种渲染格式
type User struct {
	Name  string `json:"name" xml:"name" yaml:"name"`
	Email string `json:"email" xml:"email" yaml:"email"`
	Age   int    `json:"age" xml:"age" yaml:"age"`
}

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 加载 HTML 模板
	router.LoadHTMLGlob("templates/*")

	// 1. JSON 渲染
	router.GET("/json", func(c *gin.Context) {
		user := User{
			Name:  "Alice",
			Email: "alice@example.com",
			Age:   25,
		}

		// 使用 c.JSON() 渲染 JSON
		c.JSON(http.StatusOK, user)
	})

	// 2. 结构化 JSON（使用 gin.H）
	router.GET("/json-map", func(c *gin.Context) {
		// gin.H 是 map[string]interface{} 的快捷方式
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello",
			"status":  "success",
			"data": gin.H{
				"user": "Alice",
				"role": "admin",
			},
		})
	})

	// 3. XML 渲染
	router.GET("/xml", func(c *gin.Context) {
		user := User{
			Name:  "Bob",
			Email: "bob@example.com",
			Age:   30,
		}

		// 使用 c.XML() 渲染 XML
		c.XML(http.StatusOK, user)
	})

	// 4. YAML 渲染
	router.GET("/yaml", func(c *gin.Context) {
		user := User{
			Name:  "Charlie",
			Email: "charlie@example.com",
			Age:   35,
		}

		// 使用 c.YAML() 渲染 YAML
		c.YAML(http.StatusOK, user)
	})

	// 5. ProtoBuf 渲染
	// 注意：需要先定义 .proto 文件并生成 Go 代码
	// 这里用 map 模拟
	router.GET("/protobuf", func(c *gin.Context) {
		// 实际使用时应该使用 protobuf 生成的结构体
		data := map[string]interface{}{
			"name":  "David",
			"email": "david@example.com",
			"age":   40,
		}

		// 注意：这里为了演示使用 JSON，实际应使用 c.ProtoBuf()
		c.JSON(http.StatusOK, gin.H{
			"note": "实际使用时应该使用 c.ProtoBuf()",
			"data": data,
		})
	})

	// 6. String 渲染（纯文本）
	router.GET("/string", func(c *gin.Context) {
		// 使用 c.String() 渲染纯文本
		c.String(http.StatusOK, "Hello, %s! You are %d years old.", "Alice", 25)
	})

	// 7. HTML 渲染
	router.GET("/html", func(c *gin.Context) {
		// 使用 c.HTML() 渲染 HTML 模板
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "Gin Rendering",
			"message": "Welcome to Gin!",
		})
	})

	// 8. Data 渲染（原始二进制数据）
	router.GET("/data", func(c *gin.Context) {
		// 使用 c.Data() 渲染原始二进制数据
		data := []byte("This is raw binary data")
		c.Data(http.StatusOK, "application/octet-stream", data)
	})

	// 9. File 渲染（文件下载）
	router.GET("/file", func(c *gin.Context) {
		// 使用 c.File() 发送文件
		// 注意：这里使用当前文件作为示例
		c.File("main.go")
	})

	// 10. FileAttachment（强制下载）
	router.GET("/download", func(c *gin.Context) {
		// 使用 c.FileAttachment() 强制浏览器下载文件
		c.FileAttachment("main.go", "gin-example.go")
	})

	// 11. 根据 Accept 头自动选择格式
	router.GET("/auto", func(c *gin.Context) {
		user := User{
			Name:  "Eve",
			Email: "eve@example.com",
			Age:   28,
		}

		// 根据客户端的 Accept 头返回不同格式
		accept := c.GetHeader("Accept")

		switch accept {
		case "application/xml":
			c.XML(http.StatusOK, user)
		case "application/x-yaml":
			c.YAML(http.StatusOK, user)
		default:
			c.JSON(http.StatusOK, user)
		}
	})

	// 12. 自定义响应头
	router.GET("/custom-headers", func(c *gin.Context) {
		// 设置自定义响应头
		c.Header("X-Custom-Header", "CustomValue")
		c.Header("X-Request-ID", "12345")

		c.JSON(http.StatusOK, gin.H{
			"message": "Check the response headers",
		})
	})

	// 13. 设置 Content-Type
	router.GET("/custom-content-type", func(c *gin.Context) {
		// 手动设置 Content-Type
		c.Header("Content-Type", "application/vnd.api+json")

		c.String(http.StatusOK, `{"data": {"type": "users", "id": "1"}}`)
	})

	// 14. 流式响应
	router.GET("/stream", func(c *gin.Context) {
		// 设置响应头
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		// 模拟流式数据
		for i := 0; i < 5; i++ {
			c.SSEvent("message", gin.H{
				"time": i,
				"text": "Hello",
			})
			c.Writer.Flush()
		}
	})

	// 15. 条件渲染
	router.GET("/conditional", func(c *gin.Context) {
		format := c.DefaultQuery("format", "json")

		data := gin.H{
			"message": "Conditional rendering",
			"format":  format,
		}

		switch format {
		case "xml":
			c.XML(http.StatusOK, data)
		case "yaml":
			c.YAML(http.StatusOK, data)
		case "text":
			c.String(http.StatusOK, "Message: %s, Format: %s", data["message"], data["format"])
		default:
			c.JSON(http.StatusOK, data)
		}
	})

	// 打印测试说明
	println("=========================================")
	println("渲染示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 支持多种数据格式：JSON, XML, YAML, ProtoBuf")
	println("- 支持文本、HTML、文件渲染")
	println("- 支持自定义响应头和 Content-Type")
	println("- 支持流式响应和条件渲染")
	println("")
	println("测试命令：")
	println("")
	println("1. JSON 渲染：")
	println("   curl http://localhost:8080/json")
	println("")
	println("2. XML 渲染：")
	println("   curl http://localhost:8080/xml")
	println("")
	println("3. YAML 渲染：")
	println("   curl http://localhost:8080/yaml")
	println("")
	println("4. 纯文本渲染：")
	println("   curl http://localhost:8080/string")
	println("")
	println("5. HTML 渲染：")
	println("   curl http://localhost:8080/html")
	println("")
	println("6. 文件下载：")
	println("   curl http://localhost:8080/file")
	println("   curl http://localhost:8080/download")
	println("")
	println("7. 根据 Accept 头自动选择：")
	println("   curl -H 'Accept: application/json' http://localhost:8080/auto")
	println("   curl -H 'Accept: application/xml' http://localhost:8080/auto")
	println("")
	println("8. 条件渲染：")
	println("   curl 'http://localhost:8080/conditional?format=json'")
	println("   curl 'http://localhost:8080/conditional?format=xml'")
	println("   curl 'http://localhost:8080/conditional?format=yaml'")
	println("")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
