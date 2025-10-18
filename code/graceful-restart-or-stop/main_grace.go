package main

import (
	"log"
	"net/http"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 定义路由
	router.GET("/", func(c *gin.Context) {
		time.Sleep(2 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server with Grace")
	})

	// 打印启动信息
	log.Println("========================================")
	log.Println("使用 grace 实现优雅重启")
	log.Println("========================================")
	log.Println("服务器启动在 :8080")
	log.Println("")
	log.Println("功能说明：")
	log.Println("- grace 是 Facebook 开源的优雅重启库")
	log.Println("- 支持零停机时间的重启")
	log.Println("- 可以同时监听多个端口")
	log.Println("")
	log.Println("重要提示：")
	log.Println("- grace 不支持 'go run' 模式")
	log.Println("- 必须先编译成可执行文件再运行")
	log.Println("")
	log.Println("正确的测试步骤：")
	log.Println("1. 编译：go build -o server main_grace.go")
	log.Println("2. 运行：./server")
	log.Println("3. 测试请求：curl http://localhost:8080/")
	log.Println("4. 查看进程：ps aux | grep server")
	log.Println("5. 发送 SIGUSR2 重启：kill -USR2 <pid>")
	log.Println("6. 或按 Ctrl+C 优雅关闭")
	log.Println("")
	log.Println("注意：")
	log.Println("- SIGUSR2: 触发优雅重启（零停机）")
	log.Println("- SIGINT/SIGTERM: 触发优雅关闭")
	log.Println("========================================")

	// 使用 gracehttp.Serve 替代标准的 ListenAndServe
	// gracehttp 会自动处理信号并实现优雅重启/关闭
	//
	// 特点：
	// - 接收到 SIGUSR2 时，fork 新进程并优雅关闭旧进程
	// - 接收到 SIGINT/SIGTERM 时，优雅关闭当前进程
	// - 支持同时监听多个端口
	if err := gracehttp.Serve(
		&http.Server{Addr: ":8080", Handler: router},
	); err != nil {
		log.Fatalf("listen: %s\n", err)
	}

	log.Println("服务器已退出")
}
