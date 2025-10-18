package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义两个不同的结构体
type formA struct {
	Foo string `json:"foo" xml:"foo" binding:"required"`
}

type formB struct {
	Bar string `json:"bar" xml:"bar" binding:"required"`
}

// 错误的处理方式：尝试多次绑定同一个 request body
func wrongHandler(c *gin.Context) {
	objA := formA{}
	objB := formB{}

	// 第一次尝试绑定到 formA
	// c.ShouldBind 使用了 c.Request.Body，不可重用
	if errA := c.ShouldBind(&objA); errA == nil {
		c.String(http.StatusOK, `the body should be formA`)
		return
	}

	// 因为现在 c.Request.Body 已经是 EOF，所以这里会报错
	// 即使发送的数据符合 formB 的结构，这里也会失败
	if errB := c.ShouldBind(&objB); errB == nil {
		c.String(http.StatusOK, `the body should be formB`)
		return
	}

	// 两次绑定都失败
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Invalid request - both bindings failed",
		"note":  "This is expected because Request.Body can only be read once",
	})
}

func main() {
	router := gin.Default()

	// 注册路由
	router.POST("/test", wrongHandler)

	// 打印测试说明
	println("========================================")
	println("错误示例：多次绑定 Request Body")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 formA 格式（会成功）：")
	println(`   curl -X POST http://localhost:8080/test \`)
	println(`     -H "Content-Type: application/json" \`)
	println(`     -d '{"foo":"hello"}'`)
	println("")
	println("2. 测试 formB 格式（会失败，因为第二次绑定无法读取 body）：")
	println(`   curl -X POST http://localhost:8080/test \`)
	println(`     -H "Content-Type: application/json" \`)
	println(`     -d '{"bar":"world"}'`)
	println("")
	println("预期结果：")
	println("- formA 格式会成功返回")
	println("- formB 格式会失败，因为第一次绑定已经消耗了 Request.Body")
	println("========================================")
	println("")

	// 启动服务器
	router.Run(":8080")
}
