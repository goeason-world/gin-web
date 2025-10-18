package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的 Gin 引擎
	router := gin.Default()

	// 定义一个 POST 路由 /post
	router.POST("/post", func(c *gin.Context) {
		// c.QueryMap() 用于解析 URL 查询字符串中的 map
		// 例如, /post?ids[a]=123&ids[b]=456
		// "ids" 是 map 的键
		ids := c.QueryMap("ids")

		// c.PostFormMap() 用于解析 POST 表单中的 map
		// "names" 是 map 的键
		names := c.PostFormMap("names")

		// 在控制台打印出解析后的 map，方便调试
		fmt.Printf("URL Query IDs: %v\n", ids)
		fmt.Printf("Post Form Names: %v\n", names)

		// 将解析后的数据以 JSON 格式返回给客户端
		c.JSON(http.StatusOK, gin.H{
			"message": "Data received successfully!",
			"query_ids":   ids,
			"form_names":  names,
		})
	})

	// 打印测试说明
	println("=========================================")
	println("将查询字符串或表单参数映射为 Map 示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("   curl --globoff -X POST 'http://localhost:8080/post?ids[a]=123&ids[b]=hello' \\")
	println("     -H 'Content-Type: application/x-www-form-urlencoded' \\")
	println("     -d 'names[first]=alice&names[second]=bob'")
	println("")
	println("预期输出：")
	println("   客户端将收到包含解析后数据的 JSON 响应。")
	println("   服务器控制台将打印出解析后的 map。")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}