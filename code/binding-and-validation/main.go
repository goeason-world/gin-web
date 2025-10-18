package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Login 结构体定义了登录请求的数据结构
// 使用各种验证标签来确保数据的有效性
type Login struct {
	// User 字段必需，且必须是有效的邮箱格式
	User string `form:"user" json:"user" xml:"user" binding:"required,email"`
	// Password 字段必需，且长度必须至少为 6 个字符
	Password string `form:"password" json:"password" xml:"password" binding:"required,min=6"`
}

// Booking 结构体展示了自定义验证器的使用
type Booking struct {
	// CheckIn 和 CheckOut 使用自定义的 bookabledate 验证器
	CheckIn  string `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut string `form:"check_out" binding:"required,gtfield=CheckIn,bookabledate" time_format:"2006-01-02"`
}

// bookableDate 是一个自定义验证函数
// 它检查日期是否在今天之后（即是否可预订）
var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	// 这里简化处理，实际应用中应该解析日期并与当前日期比较
	// 为了演示，我们只检查字符串是否非空
	date, ok := fl.Field().Interface().(string)
	if ok {
		// 实际应用中，这里应该解析日期并验证
		// 例如：time.Parse("2006-01-02", date)
		// 然后检查是否在今天之后
		return len(date) > 0
	}
	return false
}

func main() {
	route := gin.Default()

	// 注册自定义验证器
	// 将 "bookabledate" 标签与 bookableDate 函数关联
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}

	// 登录路由 - 演示基本的绑定和验证
	route.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		// ShouldBindJSON 只绑定 JSON 格式的数据
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 验证通过，检查用户名和密码（这里简化处理）
		if json.User != "test@example.com" || json.Password != "password123" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "you are logged in",
		})
	})

	// 登录路由 - 支持 XML 格式
	route.POST("/loginXML", func(c *gin.Context) {
		var xml Login
		// ShouldBindXML 只绑定 XML 格式的数据
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if xml.User != "test@example.com" || xml.Password != "password123" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "you are logged in",
		})
	})

	// 登录路由 - 支持表单格式
	route.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// ShouldBind 会根据 Content-Type 自动选择绑定方式
		// 对于表单数据，它会绑定 form-urlencoded 或 multipart/form-data
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if form.User != "test@example.com" || form.Password != "password123" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "you are logged in",
		})
	})

	// 预订路由 - 演示自定义验证器
	route.GET("/bookable", func(c *gin.Context) {
		var b Booking
		// 从查询参数中绑定数据
		if err := c.ShouldBindWith(&b, binding.Query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Booking dates are valid!",
			"check_in":  b.CheckIn,
			"check_out": b.CheckOut,
		})
	})

	// 打印测试说明
	println("========================================")
	println("模型绑定和验证示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 JSON 登录（成功）：")
	println(`   curl -X POST http://localhost:8080/loginJSON \`)
	println(`     -H "Content-Type: application/json" \`)
	println(`     -d '{"user":"test@example.com","password":"password123"}'`)
	println("")
	println("2. 测试 JSON 登录（邮箱格式错误）：")
	println(`   curl -X POST http://localhost:8080/loginJSON \`)
	println(`     -H "Content-Type: application/json" \`)
	println(`     -d '{"user":"invalid-email","password":"password123"}'`)
	println("")
	println("3. 测试 JSON 登录（密码太短）：")
	println(`   curl -X POST http://localhost:8080/loginJSON \`)
	println(`     -H "Content-Type: application/json" \`)
	println(`     -d '{"user":"test@example.com","password":"123"}'`)
	println("")
	println("4. 测试 XML 登录：")
	println(`   curl -X POST http://localhost:8080/loginXML \`)
	println(`     -H "Content-Type: application/xml" \`)
	println(`     -d '<login><user>test@example.com</user><password>password123</password></login>'`)
	println("")
	println("5. 测试表单登录：")
	println(`   curl -X POST http://localhost:8080/loginForm \`)
	println(`     -d "user=test@example.com&password=password123"`)
	println("")
	println("6. 测试预订日期验证：")
	println(`   curl "http://localhost:8080/bookable?check_in=2025-01-15&check_out=2025-01-20"`)
	println("")
	println("========================================")
	println("")

	// 监听 8088 端口
	route.Run(":8080")
}
