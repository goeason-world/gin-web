package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// auditEvent 审计事件结构体
// 捕获我们想要异步记录的最小信息集
type auditEvent struct {
	Path     string        // 请求路径，如 /api/users
	Method   string        // HTTP 方法，如 GET、POST
	Status   int           // HTTP 状态码，如 200、404
	ClientIP string        // 客户端 IP 地址
	Latency  time.Duration // 请求处理耗时
	When     time.Time     // 事件发生时间
}

// auditLogger 审计日志器
// 负责将审计事件排队并在 worker goroutine 中处理
//
// 设计思路：
// 1. 使用带缓冲的 channel 作为队列，避免阻塞请求
// 2. 使用多个 worker goroutine 并发处理日志
// 3. 使用 WaitGroup 确保优雅退出时所有日志都被处理
type auditLogger struct {
	queue     chan auditEvent // 事件队列（带缓冲的 channel）
	wg        sync.WaitGroup  // 用于等待所有 worker 完成
	closeOnce sync.Once       // 确保 Close 只执行一次
}

// newAuditLogger 创建一个后台审计日志器
//
// 参数：
//
//	buffer: 队列缓冲区大小，决定了可以暂存多少个待处理的事件
//	workers: worker goroutine 的数量，决定了并发处理能力
//
// 返回：
//
//	*auditLogger: 已启动的审计日志器实例
func newAuditLogger(buffer int, workers int) *auditLogger {
	// 步骤1：创建审计日志器实例
	logger := &auditLogger{
		queue: make(chan auditEvent, buffer), // 创建带缓冲的 channel
	}

	// 步骤2：启动多个 worker goroutine 来处理日志
	// 多个 worker 可以并发处理，提高吞吐量
	for i := 0; i < workers; i++ {
		logger.wg.Add(1) // 增加 WaitGroup 计数

		// 启动 worker goroutine
		go func(id int) {
			defer logger.wg.Done() // goroutine 结束时减少计数

			// 持续从队列中读取事件并处理
			// 当 channel 关闭且为空时，range 循环会自动退出
			for event := range logger.queue {
				// 实际的日志处理逻辑
				// 在生产环境中，这里可以：
				// - 写入数据库
				// - 发送到日志收集系统（如 ELK）
				// - 推送到消息队列
				log.Printf(
					"[audit-worker-%d] %s %s status=%d latency=%s ip=%s",
					id,
					event.Method,
					event.Path,
					event.Status,
					event.Latency,
					event.ClientIP,
				)
			}

			// 队列关闭后，worker 退出
			log.Printf("[audit-worker-%d] queue closed", id)
		}(i + 1) // 传入 worker ID（从1开始）
	}

	return logger
}

// Middleware 返回一个 Gin 中间件函数
// 该中间件会将审计事件推送到队列中
//
// 工作流程：
// 1. 记录请求开始时间
// 2. 执行后续的处理器
// 3. 收集审计信息
// 4. 非阻塞地将事件推送到队列
func (a *auditLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 步骤1：记录请求开始时间
		// 用于计算请求处理耗时
		start := time.Now()

		// 步骤2：执行后续的处理器（包括路由处理函数）
		// c.Next() 会依次执行处理器链中的下一个处理器
		// 当所有处理器执行完毕后，才会继续执行下面的代码
		c.Next()

		// 步骤3：收集审计信息
		// 此时请求已经处理完毕，可以获取完整的请求和响应信息
		event := auditEvent{
			Path:     c.Request.URL.Path, // 请求路径
			Method:   c.Request.Method,   // HTTP 方法
			Status:   c.Writer.Status(),  // 响应状态码
			ClientIP: c.ClientIP(),       // 客户端 IP
			Latency:  time.Since(start),  // 请求处理耗时
			When:     time.Now(),         // 当前时间
		}

		// 步骤4：非阻塞地将事件推送到队列
		// 使用 select 语句实现非阻塞发送
		select {
		case a.queue <- event:
			// 成功推送到队列，不阻塞请求
			// worker goroutine 会异步处理这个事件

		default:
			// 队列已满，无法推送
			// 采用同步方式记录日志，避免数据丢失
			// 这是一个降级策略：宁可稍微影响性能，也不能丢失重要的审计日志
			log.Printf(
				"[audit] queue full: %s %s status=%d latency=%s",
				event.Method,
				event.Path,
				event.Status,
				event.Latency,
			)
		}

		// 注意：这里不等待日志处理完成，请求会立即返回
	}
}

// Close 停止所有 worker 并等待它们完成处理
//
// 优雅退出流程：
// 1. 关闭队列 channel（不再接受新事件）
// 2. 等待所有 worker 处理完队列中的剩余事件
// 3. 确保没有日志丢失
//
// 注意：使用 sync.Once 确保只执行一次，避免重复关闭 channel
func (a *auditLogger) Close() {
	a.closeOnce.Do(func() {
		// 步骤1：关闭队列 channel
		// 关闭后，不能再向 channel 发送数据
		// 但 worker 仍然可以读取 channel 中的剩余数据
		close(a.queue)

		// 步骤2：等待所有 worker goroutine 完成
		// WaitGroup.Wait() 会阻塞，直到所有 worker 都调用了 Done()
		// 这确保了所有待处理的日志都被处理完毕
		a.wg.Wait()

		// 此时所有 worker 都已退出，可以安全地关闭程序
	})
}

func main() {
	// 步骤1：设置 Gin 为 release 模式
	// release 模式下日志输出更简洁，性能更好
	gin.SetMode(gin.ReleaseMode)

	// 步骤2：创建一个新的 Gin 引擎（不包含默认中间件）
	router := gin.New()

	// 步骤3：手动添加 Recovery 中间件
	// Recovery 中间件用于捕获 panic，防止程序崩溃
	router.Use(gin.Recovery())

	// 步骤4：创建审计日志器
	// 参数说明：
	//   16: 队列缓冲区大小，可以暂存16个待处理的事件
	//   2:  启动2个 worker goroutine 并发处理日志
	auditor := newAuditLogger(16, 2)

	// 步骤5：使用 defer 确保程序退出时关闭审计日志器
	// 这样可以保证所有待处理的日志都被处理完毕
	defer auditor.Close()

	// 步骤6：安装审计日志中间件
	// 这个中间件会在每个请求处理后收集审计信息
	router.Use(auditor.Middleware())

	// 路由1：生成报告接口
	// 模拟一个中等耗时的操作（500ms）
	router.GET("/report", func(c *gin.Context) {
		// 模拟报告生成过程
		time.Sleep(500 * time.Millisecond)

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"message": "report generated",
		})
		// 请求返回后，审计日志会在后台异步处理
	})

	// 路由2：创建任务接口
	// 模拟一个较耗时的操作（1.5秒）
	router.POST("/jobs", func(c *gin.Context) {
		// 模拟任务创建过程
		time.Sleep(1500 * time.Millisecond)

		// 返回创建成功响应
		c.JSON(http.StatusCreated, gin.H{
			"message": "job accepted",
		})
	})

	// 打印使用说明
	log.Println("========================================")
	log.Println("Goroutines Inside Middleware (扩展版)")
	log.Println("========================================")
	log.Println("- 使用后台 worker 处理审计日志，演示如何管理 goroutine 生命周期")
	log.Println("- 通过信号优雅退出，确保后台任务完成")
	log.Println("")
	log.Println("测试命令：")
	log.Println("   curl http://localhost:8080/report")
	log.Println("   curl -X POST http://localhost:8080/jobs")
	log.Println("")
	log.Println("按 Ctrl+C 可以触发优雅退出流程")
	log.Println("========================================")

	// 步骤7：创建 HTTP 服务器
	srv := &http.Server{
		Addr:    ":8080", // 监听地址和端口
		Handler: router,  // 使用 Gin 路由作为处理器
	}

	// 步骤8：启动一个 goroutine 来处理优雅退出
	// 这个 goroutine 会监听系统信号（Ctrl+C 或 kill 命令）
	go func() {
		// 创建一个接收系统信号的 channel
		stop := make(chan os.Signal, 1)

		// 注册要监听的信号
		// SIGINT: Ctrl+C 触发
		// SIGTERM: kill 命令触发
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		// 阻塞等待信号
		<-stop

		// 收到退出信号，开始优雅退出流程
		log.Println("[main] shutting down...")

		// 创建一个带超时的 context
		// 给服务器5秒时间来完成正在处理的请求
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 优雅关闭服务器
		// Shutdown 会：
		// 1. 停止接受新的连接
		// 2. 等待所有活跃的连接处理完毕
		// 3. 如果超时，强制关闭
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("[main] shutdown error: %v", err)
		}

		// 注意：main 函数中的 defer auditor.Close() 会在这之后执行
		// 确保所有审计日志都被处理完毕
	}()

	// 步骤9：启动 HTTP 服务器
	// ListenAndServe 会阻塞，直到服务器关闭
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// 如果是正常关闭（http.ErrServerClosed），不报错
		// 其他错误则记录并退出
		log.Fatalf("server error: %v", err)
	}

	// 程序退出时的执行顺序：
	// 1. 收到退出信号
	// 2. srv.Shutdown() 优雅关闭 HTTP 服务器
	// 3. ListenAndServe() 返回
	// 4. defer auditor.Close() 执行，关闭审计日志器
	// 5. 等待所有 worker 处理完剩余日志
	// 6. 程序完全退出
}
