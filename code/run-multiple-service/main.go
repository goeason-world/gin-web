package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// 全局错误组，用于管理多个 goroutine 的生命周期
// errgroup.Group 提供了一种优雅的方式来等待一组 goroutine 完成
// 并且会在任何一个 goroutine 返回错误时取消其他所有 goroutine
var (
	g errgroup.Group
)

// router01 创建第一个服务的路由处理器
// 返回 http.Handler 接口，可以直接用于 http.Server
func router01() http.Handler {
	// 创建一个新的 Gin 引擎实例，不包含任何默认中间件
	// 这样可以更精确地控制每个服务的中间件配置
	e := gin.New()

	// 添加 Recovery 中间件，用于捕获 panic 并恢复程序运行
	// 这是生产环境中必须的中间件，防止单个请求的 panic 导致整个服务崩溃
	e.Use(gin.Recovery())

	// 定义根路径的 GET 请求处理器
	e.GET("/", func(c *gin.Context) {
		// 返回 JSON 响应，标识这是第一个服务
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":    http.StatusOK,
				"message": "Welcome to server 01",
				"service": "service-01",
				"port":    "8080",
			},
		)
	})

	// 添加一个健康检查端点
	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "service-01",
		})
	})

	return e
}

// router02 创建第二个服务的路由处理器
// 与 router01 类似，但提供不同的响应内容
func router02() http.Handler {
	// 创建第二个独立的 Gin 引擎实例
	e := gin.New()

	// 添加 Recovery 中间件
	e.Use(gin.Recovery())

	// 定义根路径的 GET 请求处理器
	e.GET("/", func(c *gin.Context) {
		// 返回 JSON 响应，标识这是第二个服务
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":    http.StatusOK,
				"message": "Welcome to server 02",
				"service": "service-02",
				"port":    "8081",
			},
		)
	})

	// 添加一个健康检查端点
	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "service-02",
		})
	})

	return e
}

func main() {
	// 创建第一个 HTTP 服务器实例
	// 配置监听地址、处理器和超时参数
	server01 := &http.Server{
		Addr:         ":8080",          // 监听端口 8080
		Handler:      router01(),       // 使用 router01 作为请求处理器
		ReadTimeout:  5 * time.Second,  // 读取请求的超时时间
		WriteTimeout: 10 * time.Second, // 写入响应的超时时间
	}

	// 创建第二个 HTTP 服务器实例
	// 配置不同的端口，但使用相同的超时设置
	server02 := &http.Server{
		Addr:         ":8081",          // 监听端口 8081
		Handler:      router02(),       // 使用 router02 作为请求处理器
		ReadTimeout:  5 * time.Second,  // 读取请求的超时时间
		WriteTimeout: 10 * time.Second, // 写入响应的超时时间
	}

	// 使用 errgroup 启动第一个服务器
	// g.Go() 会在新的 goroutine 中执行传入的函数
	// 如果函数返回错误，errgroup 会记录这个错误
	g.Go(func() error {
		log.Println("Starting server 01 on :8080")
		// ListenAndServe 会阻塞当前 goroutine 直到服务器关闭或出现错误
		return server01.ListenAndServe()
	})

	// 使用 errgroup 启动第二个服务器
	g.Go(func() error {
		log.Println("Starting server 02 on :8081")
		// 同时启动第二个服务器
		return server02.ListenAndServe()
	})

	// 等待所有 goroutine 完成
	// Wait() 会阻塞直到所有通过 g.Go() 启动的 goroutine 都返回
	// 如果任何一个 goroutine 返回错误，Wait() 会立即返回第一个遇到的错误
	if err := g.Wait(); err != nil {
		// 如果有任何服务器启动失败，记录错误并退出程序
		log.Fatal("Failed to start servers:", err)
	}
}
