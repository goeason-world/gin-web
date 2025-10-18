package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 定义两个不同的结构体
type formA struct {
	Foo string `json:"foo" xml:"foo" binding:"required"`
}

type formB struct {
	Bar string `json:"bar" xml:"bar" binding:"required"`
}

// 正确的处理方式：使用 ShouldBindBodyWith 允许多次绑定
func correctHandler(c *gin.Context) {
	objA := formA{}
	objB := formB{}

	// 使用 ShouldBindBodyWith 读取 c.Request.Body 并将结果存入上下文
	// 第一个参数是目标结构体，第二个参数指定绑定类型
	if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "the body should be formA",
			"data":    objA,
		})
		return
	}

	// 这时，复用存储在上下文中的 body
	// 即使第一次绑定失败，第二次仍然可以读取到完整的 body
	if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "the body should be formB JSON",
			"data":    objB,
		})
		return
	}

	// 可以接受其他格式，比如 XML
	if errB2 := c.ShouldBindBodyWith(&objB, binding.XML); errB2 == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "the body should be formB XML",
			"data":    objB,
		})
		return
	}

	// 所有绑定都失败
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Invalid request format",
		"note":  "Body doesn't match formA or formB structure",
	})
}

func main() {
	router := gin.Default()

	// 注册路由
	router.POST("/test", correctHandler)

	// 打印测试说明
	println("========================================")
	println("正确示例：使用 ShouldBindBodyWith 多次绑定")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 formA 格式（JSON）：")
	println(`   curl -X POST http://localhost:8080/test \`)
	println(`     -H "Content-Type: application/json" \`)
	println(`     -d '{"foo":"hello"}'`)
	println("")
	println("2. 测试 formB 格式（JSON）：")
	println(`   curl -X POST http://localhost:8080/test \`)
	println(`     -H "Content-Type: application/json" \`)
	println(`     -d '{"bar":"world"}'`)
	println("")
	println("3. 测试 formB 格式（XML）：")
	println(`   curl -X POST http://localhost:8080/test \`)
	println(`     -H "Content-Type: application/xml" \`)
	println(`     -d '<formB><bar>world</bar></formB>'`)
	println("")
	println("4. 测试无效格式：")
	println(`   curl -X POST http://localhost:8080/test \`)
	println(`     -H "Content-Type: application/json" \`)
	println(`     -d '{"invalid":"data"}'`)
	println("")
	println("预期结果：")
	println("- 所有格式都能正确识别和绑定")
	println("- 即使第一次绑定失败，后续绑定仍然可以成功")
	println("========================================")
	println("")

	// 启动服务器
	router.Run(":8080")
}
