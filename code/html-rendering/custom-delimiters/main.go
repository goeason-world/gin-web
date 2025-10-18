package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Delims("{{", "}}")
	router.LoadHTMLGlob("templates/*")

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Custom Delimiters",
		})
	})

	println("========================================")
	println("自定义分隔符示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("   curl http://localhost:8080/index")
	println("")

	router.Run(":8080")
}
