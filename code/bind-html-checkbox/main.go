package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MyForm 定义了我们表单数据的结构
// `form:"colors[]"` 这个标签是关键，它告诉 Gin 将所有名为 "colors[]" 的
// 表单值收集到一个切片中。
type MyForm struct {
	Colors []string `form:"colors[]"`
}

func main() {
	router := gin.Default()

	// 从 "templates" 目录加载 HTML 模板
	router.LoadHTMLGlob("templates/*")

	// GET 路由，用于显示表单页面
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form.html", nil)
	})

	// POST 路由，用于处理表单提交
	router.POST("/", func(c *gin.Context) {
		var form MyForm
		// 将表单数据绑定到 MyForm 结构体
		// c.ShouldBind() 会自动处理 form-urlencoded 类型的数据
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 为演示目的，在控制台打印出所选的颜色
		fmt.Printf("Selected colors: %v\n", form.Colors)

		// 将所选的颜色作为 JSON 响应返回
		c.JSON(http.StatusOK, gin.H{
			"message": "Form submitted successfully!",
			"colors":  form.Colors,
		})
	})

	// 在 8080 端口启动服务
	router.Run(":8080")
}
