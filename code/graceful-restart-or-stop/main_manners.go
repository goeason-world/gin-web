package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manners/manners"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 定义路由
	router.GET("/", func(c *gin.Context) {
		time.Sleep(2 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server with Manners")
	})

	// 打印启动信息
	log.Println("========================================")
	log.Println("使用 manners 实现优雅关闭")
	log.Println("========================================")
	log.Println("服务器启动在 :8080")
	log.Println("")
	log.Println("功能说明：")
	log.Println("- manners 提供了优雅关闭的能力")
	log.Println("- 会等待所有活跃连接完成后再关闭")
	log.Println("")
	log.Println("测试步骤：")
	log.Println("1. 在新终端执行：curl http://localhost:8080/")
	log.Println("2. 在请求处理过程中按 Ctrl+C")
	log.Println("3. 观察服务器会等待请求完成后再退出")
	log.Println("========================================")

	// 在一个 goroutine 中启动服务器
	go func() {
		// 使用 manners.ListenAndServe 替代标准的 http.ListenAndServe
		// manners 会跟踪所有活跃的连接
		if err := manners.ListenAndServe(":8080", router); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 使用 manners.Close() 优雅关闭服务器
	// Close 会：
	// 1. 停止接受新连接
	// 2. 等待所有活跃连接完成
	// 3. 没有超时限制，会一直等待
	manners.Close()

	log.Println("服务器已退出")
}
