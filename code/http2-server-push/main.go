package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var html = template.Must(template.New("https").Parse(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HTTP/2 Server Push 示例</title>
    <link rel="stylesheet" href="/assets/app.css">
</head>
<body>
    <div class="container">
        <h1>HTTP/2 Server Push 示例</h1>
        <p>此页面演示了 HTTP/2 Server Push 功能。</p>
        <p>服务器会主动推送 CSS 和 JS 文件，而不是等待浏览器请求。</p>
        <div class="info">
            <h2>什么是 HTTP/2 Server Push？</h2>
            <p>HTTP/2 Server Push 允许服务器在客户端请求之前主动推送资源，减少加载时间。</p>
        </div>
        <img src="/assets/logo.png" alt="Logo" class="logo">
    </div>
    <script src="/assets/app.js"></script>
</body>
</html>
`))

func main() {
	// 创建默认的 Gin 引擎
	r := gin.Default()

	// 静态文件服务
	r.Static("/assets", "./assets")

	// 主页路由，使用 HTTP/2 Server Push
	r.GET("/", func(c *gin.Context) {
		// 检查是否支持 HTTP/2 Server Push
		// Pusher 接口只在 HTTP/2 连接中可用
		if pusher := c.Writer.Pusher(); pusher != nil {
			// 推送 CSS 文件
			// Server Push 的好处是浏览器无需解析 HTML 就能提前获取这些资源
			if err := pusher.Push("/assets/app.css", nil); err != nil {
				log.Printf("Failed to push app.css: %v", err)
			}

			// 推送 JavaScript 文件
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push app.js: %v", err)
			}

			// 推送图片文件
			// 可以推送任何类型的资源
			if err := pusher.Push("/assets/logo.png", nil); err != nil {
				log.Printf("Failed to push logo.png: %v", err)
			}

			log.Println("HTTP/2 Server Push: 已推送 CSS、JS 和图片资源")
		} else {
			// 如果不支持 HTTP/2，会降级为普通的 HTTP/1.1
			log.Println("HTTP/2 Server Push 不可用，使用 HTTP/1.1")
		}

		// 渲染 HTML 模板
		if err := html.Execute(c.Writer, nil); err != nil {
			log.Printf("Failed to execute template: %v", err)
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	})

	// 打印测试说明
	println("========================================")
	println("HTTP/2 Server Push 示例")
	println("========================================")
	println("服务器启动在 https://localhost:8080")
	println("")
	println("功能说明：")
	println("- 演示 HTTP/2 Server Push 功能")
	println("- 服务器主动推送 CSS、JS 和图片资源")
	println("- 减少页面加载时间，提升用户体验")
	println("")
	println("测试步骤：")
	println("")
	println("1. 生成自签名证书（仅用于测试）：")
	println("   openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 \\")
	println("     -keyout server.key -out server.crt \\")
	println("     -subj \"/CN=localhost\"")
	println("")
	println("2. 创建测试资源文件：")
	println("   mkdir -p assets")
	println("   echo 'body { font-family: Arial; }' > assets/app.css")
	println("   echo 'console.log(\"Hello HTTP/2\");' > assets/app.js")
	println("   echo 'PNG' > assets/logo.png")
	println("")
	println("3. 启动服务器：")
	println("   go run main.go")
	println("")
	println("4. 在浏览器中访问（需要接受自签名证书）：")
	println("   https://localhost:8080")
	println("")
	println("5. 打开浏览器开发者工具（F12）查看网络请求：")
	println("   - 在 Protocol 列中可以看到 \"h2\" 表示 HTTP/2")
	println("   - 推送的资源会显示 \"Push\" 标记")
	println("")
	println("注意：HTTP/2 Server Push 需要 HTTPS 连接")
	println("========================================")
	println("")

	// 使用 RunTLS 启动 HTTPS 服务器
	// HTTP/2 Server Push 需要 HTTPS 连接
	// 参数: 地址, 证书文件路径, 私钥文件路径
	if err := r.RunTLS(":8080", "./server.crt", "./server.key"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
