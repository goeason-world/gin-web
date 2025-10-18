package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// GET 请求示例
	// 用于获取资源，通常用于查询数据
	router.GET("/someGet", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "GET 请求成功",
			"method":  "GET",
		})
	})

	// POST 请求示例
	// 用于创建新资源，通常用于提交表单或上传数据
	router.POST("/somePost", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "POST 请求成功",
			"method":  "POST",
		})
	})

	// PUT 请求示例
	// 用于完整更新资源，通常需要提供完整的资源数据
	router.PUT("/somePut", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "PUT 请求成功",
			"method":  "PUT",
		})
	})

	// DELETE 请求示例
	// 用于删除资源
	router.DELETE("/someDelete", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "DELETE 请求成功",
			"method":  "DELETE",
		})
	})

	// PATCH 请求示例
	// 用于部分更新资源，只需要提供需要修改的字段
	router.PATCH("/somePatch", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "PATCH 请求成功",
			"method":  "PATCH",
		})
	})

	// HEAD 请求示例
	// 类似于 GET，但只返回响应头，不返回响应体
	// 通常用于检查资源是否存在或获取元数据
	router.HEAD("/someHead", func(c *gin.Context) {
		// HEAD 请求不应该有响应体，只设置响应头
		c.Status(http.StatusOK)
	})

	// OPTIONS 请求示例
	// 用于获取服务器支持的 HTTP 方法
	// 通常用于 CORS 预检请求
	router.OPTIONS("/someOptions", func(c *gin.Context) {
		c.Header("Allow", "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")
		c.Status(http.StatusOK)
	})

	// 打印测试说明
	println("========================================")
	println("HTTP Method 示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 演示 Gin 支持的所有标准 HTTP 方法")
	println("- 包括 GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")
	println("")
	println("测试命令：")
	println("")
	println("1. GET 请求：")
	println("   curl http://localhost:8080/someGet")
	println("")
	println("2. POST 请求：")
	println("   curl -X POST http://localhost:8080/somePost")
	println("")
	println("3. PUT 请求：")
	println("   curl -X PUT http://localhost:8080/somePut")
	println("")
	println("4. DELETE 请求：")
	println("   curl -X DELETE http://localhost:8080/someDelete")
	println("")
	println("5. PATCH 请求：")
	println("   curl -X PATCH http://localhost:8080/somePatch")
	println("")
	println("6. HEAD 请求（只返回响应头）：")
	println("   curl -I http://localhost:8080/someHead")
	println("")
	println("7. OPTIONS 请求（查看支持的方法）：")
	println("   curl -X OPTIONS -i http://localhost:8080/someOptions")
	println("")
	println("========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
