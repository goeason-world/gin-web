package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 定义一个带有单个路径参数的路由
	// :name 是路径参数，可以匹配任何非空字符串
	router.GET("/user/:name", func(c *gin.Context) {
		// 使用 c.Param() 获取路径参数的值
		name := c.Param("name")

		c.JSON(http.StatusOK, gin.H{
			"message": "Hello",
			"name":    name,
		})
	})

	// 定义一个带有多个路径参数的路由
	// :name 和 :action 都是路径参数
	router.GET("/user/:name/:action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")

		message := name + " is " + action

		c.JSON(http.StatusOK, gin.H{
			"message": message,
			"name":    name,
			"action":  action,
		})
	})

	// 定义一个带有通配符的路由
	// *action 可以匹配任意路径，包括多级路径
	// 例如: /user/john/send/email 或 /user/john/a/b/c
	router.POST("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")

		// 注意：通配符参数会包含前导斜杠
		// 例如，对于 /user/john/send，action 的值是 "/send"

		c.JSON(http.StatusOK, gin.H{
			"message": "POST request",
			"name":    name,
			"action":  action,
		})
	})

	// RESTful API 示例：用户资源的 CRUD 操作

	// 获取用户列表
	router.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "获取用户列表",
			"users": []gin.H{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
		})
	})

	// 获取单个用户
	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "获取用户详情",
			"id":      id,
			"name":    "User " + id,
		})
	})

	// 创建用户
	router.POST("/users", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"message": "创建用户成功",
			"id":      123,
		})
	})

	// 更新用户
	router.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "更新用户成功",
			"id":      id,
		})
	})

	// 删除用户
	router.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "删除用户成功",
			"id":      id,
		})
	})

	// 打印测试说明
	println("=========================================")
	println("路径参数示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 使用 :name 定义路径参数")
	println("- 使用 c.Param() 获取路径参数的值")
	println("- 支持多个路径参数")
	println("- 支持通配符路径参数")
	println("")
	println("测试命令：")
	println("")
	println("1. 单个路径参数：")
	println("   curl http://localhost:8080/user/alice")
	println("")
	println("2. 多个路径参数：")
	println("   curl http://localhost:8080/user/bob/running")
	println("")
	println("3. 通配符路径参数：")
	println("   curl -X POST http://localhost:8080/user/charlie/send/email")
	println("")
	println("4. RESTful API 示例：")
	println("   curl http://localhost:8080/users")
	println("   curl http://localhost:8080/users/123")
	println("   curl -X POST http://localhost:8080/users")
	println("   curl -X PUT http://localhost:8080/users/123")
	println("   curl -X DELETE http://localhost:8080/users/123")
	println("")
	println("预期输出：")
	println("- 每个请求都会返回包含路径参数的 JSON 响应")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
