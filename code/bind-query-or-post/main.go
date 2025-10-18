package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Person 结构体定义了我们要绑定的数据结构
// 注意 `time_format` 标签，它告诉 Gin 如何解析日期字符串
// `time_utc` 标签指定是否将时间解析为 UTC 时区
type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func main() {
	route := gin.Default()

	// 为 /testing 路由同时注册 GET 和 POST 方法
	// 两种方法都使用同一个处理函数 startPage
	// 这展示了 Gin 绑定功能的灵活性
	route.GET("/testing", startPage)
	route.POST("/testing", startPage)

	// 打印测试说明
	println("========================================")
	println("绑定查询字符串或表单数据示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 GET 请求（查询参数）：")
	println(`   curl -X GET "http://localhost:8080/testing?name=张三&address=北京市朝阳区&birthday=1990-01-15"`)
	println("")
	println("2. 测试 POST 请求（表单数据）：")
	println(`   curl -X POST "http://localhost:8080/testing" \`)
	println(`     -d "name=李四&address=上海市浦东新区&birthday=1985-06-20"`)
	println("")
	println("========================================")
	println("")

	// 监听 8085 端口
	route.Run(":8080")
}

// startPage 是一个通用的处理函数，可以同时处理 GET 和 POST 请求
func startPage(c *gin.Context) {
	var person Person

	// c.ShouldBind() 是一个智能的绑定方法，它会根据请求类型自动选择绑定方式：
	// - 如果是 GET 请求，它会绑定 URL 查询参数 (query string)
	// - 如果是 POST 请求，它会根据 Content-Type 绑定请求体数据
	//   (如 application/x-www-form-urlencoded 或 multipart/form-data)
	if err := c.ShouldBind(&person); err == nil {
		log.Println("====== 绑定成功 ======")
		log.Println("姓名:", person.Name)
		log.Println("地址:", person.Address)
		log.Println("生日:", person.Birthday)
		log.Println("===========================")
		c.String(http.StatusOK, "Success")
	} else {
		log.Println("====== 绑定失败 ======")
		log.Println(err)
		log.Println("===========================")
		c.String(http.StatusBadRequest, "Failed to bind: "+err.Error())
	}
}
