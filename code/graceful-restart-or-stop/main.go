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
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 定义一个简单的路由
	router.GET("/", func(c *gin.Context) {
		// 模拟一个需要一些时间处理的请求
		time.Sleep(2 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// 在一个 goroutine 中启动服务器
	// 这样它就不会阻塞下面的优雅关闭处理逻辑
	go func() {
		// 启动服务器
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 打印启动信息
	log.Println("========================================")
	log.Println("优雅关闭示例")
	log.Println("========================================")
	log.Println("服务器启动在 :8080")
	log.Println("")
	log.Println("测试步骤：")
	log.Println("1. 在新终端执行：curl http://localhost:8080/")
	log.Println("2. 在请求处理过程中按 Ctrl+C")
	log.Println("3. 观察服务器会等待请求完成后再退出")
	log.Println("========================================")

	// 等待中断信号以优雅地关闭服务器
	// 设置一个接收系统信号的 channel
	quit := make(chan os.Signal, 1)

	// signal.Notify 会将接收到的信号转发到 quit channel
	// 这里监听 SIGINT (Ctrl+C) 和 SIGTERM (kill 命令)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号
	<-quit
	log.Println("正在关闭服务器...")

	// 创建一个 5 秒超时的 context
	// 这意味着服务器最多等待 5 秒来完成正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	// Shutdown 会：
	// 1. 立即停止接受新的连接
	// 2. 等待所有活跃的连接处理完毕
	// 3. 如果在超时时间内完成，返回 nil
	// 4. 如果超时，返回 context.DeadlineExceeded 错误
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器已退出")
}
