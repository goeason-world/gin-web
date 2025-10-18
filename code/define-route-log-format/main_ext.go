package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 更复杂的路由日志格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		// 根据 HTTP 方法使用不同的颜色（ANSI 颜色代码）
		var methodColor string
		switch httpMethod {
		case "GET":
			methodColor = "\033[34m" // 蓝色
		case "POST":
			methodColor = "\033[32m" // 绿色
		case "PUT":
			methodColor = "\033[33m" // 黄色
		case "DELETE":
			methodColor = "\033[31m" // 红色
		default:
			methodColor = "\033[37m" // 白色
		}
		resetColor := "\033[0m"

		// 添加时间戳
		timestamp := time.Now().Format("15:04:05")

		// 输出格式化的日志
		log.Printf("[%s] %s%-7s%s | %-30s | %s (%d handlers)\n",
			timestamp,
			methodColor,
			httpMethod,
			resetColor,
			absolutePath,
			handlerName,
			nuHandlers,
		)
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome!")
	})

	router.POST("/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.PUT("/update/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"updated": true})
	})

	router.DELETE("/delete/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"deleted": true})
	})

	log.Println("========================================")
	log.Println("复杂路由日志格式示例")
	log.Println("========================================")
	log.Println("服务器启动在 :8080")
	log.Println("注意观察上面带颜色和时间戳的路由日志")
	log.Println("========================================")

	router.Run(":8080")
}
