package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Person 结构体定义了要从 URI 中绑定的字段
// 使用 `uri` 标签来指定 URI 参数的名称
type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func main() {
	route := gin.Default()

	// 定义一个带有 URI 参数的路由
	// :name 和 :id 是 URI 参数，会被绑定到 Person 结构体中
	route.GET("/:name/:id", func(c *gin.Context) {
		var person Person

		// ShouldBindUri 专门用于绑定 URI 路径参数
		// 它会从路由路径中提取参数并填充到结构体中
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 绑定成功，返回解析后的数据
		c.JSON(http.StatusOK, gin.H{
			"message": "URI 参数绑定成功！",
			"name":    person.Name,
			"id":      person.ID,
		})
	})

	// 打印测试说明
	println("========================================")
	println("绑定 URI 参数示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试有效的 UUID：")
	println(`   curl "http://localhost:8080/张三/550e8400-e29b-41d4-a716-446655440000"`)
	println("")
	println("2. 测试另一个有效的 UUID：")
	println(`   curl "http://localhost:8080/李四/6ba7b810-9dad-11d1-80b4-00c04fd430c8"`)
	println("")
	println("3. 测试无效的 UUID（会失败）：")
	println(`   curl "http://localhost:8080/王五/invalid-uuid"`)
	println("")
	println("========================================")
	println("")

	// 监听 8087 端口
	route.Run(":8080")
}
