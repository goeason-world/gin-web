package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建 Gin 路由器
	router := gin.Default()

	// 定义一些示例路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from custom HTTP server!",
		})
	})

	router.GET("/slow", func(c *gin.Context) {
		// 模拟慢速请求
		time.Sleep(5 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"message": "This was a slow request",
		})
	})

	// 自定义 HTTP 服务器配置
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,

		// 读取超时：从连接被接受到请求体完全读取的最大时间
		// 防止慢速客户端占用连接
		ReadTimeout: 10 * time.Second,

		// 写入超时：从请求头读取完成到响应写入完成的最大时间
		// 防止慢速客户端占用连接
		WriteTimeout: 10 * time.Second,

		// 空闲超时：启用 keep-alive 时，连接保持空闲状态的最大时间
		IdleTimeout: 120 * time.Second,

		// 请求头最大字节数
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 打印测试说明
	log.Println("========================================")
	log.Println("自定义 HTTP 配置示例")
	log.Println("========================================")
	log.Println("服务器配置：")
	log.Printf("- 监听地址: %s\n", srv.Addr)
	log.Printf("- 读取超时: %v\n", srv.ReadTimeout)
	log.Printf("- 写入超时: %v\n", srv.WriteTimeout)
	log.Printf("- 空闲超时: %v\n", srv.IdleTimeout)
	log.Printf("- 最大请求头: %d bytes\n", srv.MaxHeaderBytes)
	log.Println("")
	log.Println("测试命令：")
	log.Println("1. 测试正常请求：")
	log.Println("   curl http://localhost:8080/")
	log.Println("")
	log.Println("2. 测试慢速请求：")
	log.Println("   curl http://localhost:8080/slow")
	log.Println("")
	log.Println("3. 优雅关闭：按 Ctrl+C")
	log.Println("========================================")
	log.Println("")

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("服务器启动在 %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// 捕获 SIGINT (Ctrl+C) 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 创建一个 5 秒的超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器已退出")
}
