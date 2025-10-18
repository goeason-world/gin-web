package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Person 定义了查询参数的数据结构
// 使用 form 标签来指定查询参数的名称
type Person struct {
	Name    string `form:"name"`
	Address string `form:"address"`
}

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 定义一个 GET 路由，演示从查询字符串绑定数据
	router.GET("/testing", func(c *gin.Context) {
		var person Person

		// ShouldBindQuery 只绑定查询字符串参数
		// 即使请求体中有同名字段，也会被忽略
		if err := c.ShouldBindQuery(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 打印接收到的数据
		log.Printf("Name: %s, Address: %s", person.Name, person.Address)

		// 返回绑定的数据
		c.JSON(http.StatusOK, gin.H{
			"message": "查询参数绑定成功",
			"name":    person.Name,
			"address": person.Address,
		})
	})

	// 定义一个 POST 路由，演示只绑定查询字符串，忽略 POST 数据
	router.POST("/testing", func(c *gin.Context) {
		var person Person

		// ShouldBindQuery 只绑定查询字符串参数
		// POST 请求体中的数据会被忽略
		if err := c.ShouldBindQuery(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 打印接收到的数据
		log.Printf("Name: %s, Address: %s", person.Name, person.Address)

		// 返回绑定的数据
		c.JSON(http.StatusOK, gin.H{
			"message": "只绑定了查询字符串，忽略了 POST 数据",
			"name":    person.Name,
			"address": person.Address,
		})
	})

	// 打印测试说明
	println("=========================================")
	println("只绑定查询字符串示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 使用 ShouldBindQuery 只绑定查询字符串参数")
	println("- 忽略 POST 表单数据和 JSON 数据")
	println("- 确保数据只来自 URL 参数")
	println("")
	println("测试命令：")
	println("")
	println("1. GET 请求，从查询字符串获取参数：")
	println("   curl 'http://localhost:8080/testing?name=alice&address=beijing'")
	println("")
	println("2. POST 请求，只绑定查询字符串，忽略 POST 数据：")
	println("   curl -X POST 'http://localhost:8080/testing?name=bob&address=shanghai' \\")
	println("     -d 'name=ignored&address=ignored'")
	println("")
	println("3. GET 请求，缺少参数（返回空值）：")
	println("   curl 'http://localhost:8080/testing?name=charlie'")
	println("")
	println("预期输出：")
	println("- 第一个请求返回查询字符串中的值")
	println("- 第二个请求只返回查询字符串中的值，忽略 POST 数据")
	println("- 第三个请求返回部分数据，缺少的字段为空")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
