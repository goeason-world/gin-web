package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 自定义路由日志格式
	// gin.DebugPrintRouteFunc 是一个函数变量，用于自定义路由注册时的日志输出
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	router := gin.Default()

	// 注册一些示例路由
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/submit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "submitted",
		})
	})

	router.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"user_id": id,
		})
	})

	// 路由组
	api := router.Group("/api")
	{
		api.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"users": []string{"user1", "user2"},
			})
		})

		api.POST("/users", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{
				"message": "user created",
			})
		})
	}

	// 打印测试说明
	log.Println("========================================")
	log.Println("定义路由日志格式示例")
	log.Println("========================================")
	log.Println("服务器启动在 :8080")
	log.Println("")
	log.Println("注意观察上面的路由注册日志格式")
	log.Println("格式为: endpoint [方法] [路径] [处理函数] [处理器数量]")
	log.Println("")
	log.Println("测试命令：")
	log.Println("1. curl http://localhost:8080/ping")
	log.Println("2. curl -X POST http://localhost:8080/submit")
	log.Println("3. curl http://localhost:8080/user/123")
	log.Println("4. curl http://localhost:8080/api/users")
	log.Println("========================================")
	log.Println("")

	router.Run(":8080")
}
