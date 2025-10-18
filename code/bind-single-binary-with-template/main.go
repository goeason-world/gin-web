package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed assets/*
var f embed.FS

func main() {
	router := gin.Default()

	// 从嵌入的文件系统中加载模板
	// 使用 template.ParseFS 可以直接从 embed.FS 中解析模板
	templ := template.Must(template.New("").ParseFS(f, "assets/*.tmpl"))
	router.SetHTMLTemplate(templ)

	// 提供静态文件服务
	// 需要创建一个子文件系统，去掉 "assets" 前缀
	staticFS, err := fs.Sub(f, "assets")
	if err != nil {
		panic(err)
	}
	router.StaticFS("/static", http.FS(staticFS))

	// 定义路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "使用模板构建单个二进制文件",
		})
	})

	// 打印测试说明
	println("========================================")
	println("使用模板构建单个二进制文件示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 使用 Go 1.16+ 的 embed 特性")
	println("- 将模板文件嵌入到二进制文件中")
	println("- 实现单个可执行文件部署")
	println("")
	println("测试命令：")
	println("")
	println("1. 访问主页：")
	println("   curl http://localhost:8080/")
	println("")
	println("2. 或在浏览器中打开：")
	println("   http://localhost:8080/")
	println("")
	println("========================================")
	println("")

	// 监听 8086 端口
	router.Run(":8080")
}
