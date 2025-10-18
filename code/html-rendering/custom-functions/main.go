package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func main() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Custom Functions",
			"now":   time.Date(2017, 07, 04, 0, 0, 0, 0, time.UTC),
		})
	})

	println("========================================")
	println("自定义模板函数示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("   curl http://localhost:8080/")
	println("")

	router.Run(":8080")
}
