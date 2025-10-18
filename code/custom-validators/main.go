package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Booking 包含预订信息
type Booking struct {
	// CheckIn 必须是未来的日期
	CheckIn time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	// CheckOut 必须在 CheckIn 之后
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

// User 包含用户注册信息
type User struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,chinesephone"`
	Age      int    `json:"age" binding:"required,gte=18,lte=120"`
	Password string `json:"password" binding:"required,strongpassword"`
}

// bookableDate 验证日期是否在今天之后
var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		// 检查日期是否在今天之后
		if date.After(today) {
			return true
		}
	}
	return false
}

// chinesePhone 验证中国大陆手机号格式
var chinesePhone validator.Func = func(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// 简单的手机号验证：11位数字，以1开头
	if len(phone) != 11 {
		return false
	}
	if phone[0] != '1' {
		return false
	}
	// 检查是否全是数字
	for _, c := range phone {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// strongPassword 验证密码强度
// 密码必须至少8个字符，包含大写字母、小写字母和数字
var strongPassword validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func main() {
	router := gin.Default()

	// 获取验证器引擎并注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册 bookabledate 验证器
		v.RegisterValidation("bookabledate", bookableDate)
		// 注册中国手机号验证器
		v.RegisterValidation("chinesephone", chinesePhone)
		// 注册强密码验证器
		v.RegisterValidation("strongpassword", strongPassword)
	}

	// 预订路由
	router.GET("/bookings", func(c *gin.Context) {
		var b Booking
		if err := c.ShouldBindWith(&b, binding.Query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "验证失败",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "预订日期有效",
			"check_in":  b.CheckIn.Format("2006-01-02"),
			"check_out": b.CheckOut.Format("2006-01-02"),
		})
	})

	// 用户注册路由
	router.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "验证失败",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "用户注册成功",
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
			"age":      user.Age,
		})
	})

	// 打印测试说明
	fmt.Println("========================================")
	fmt.Println("自定义验证器示例")
	fmt.Println("========================================")
	fmt.Println("服务器启动在 :8080")
	fmt.Println("")
	fmt.Println("自定义验证器说明：")
	fmt.Println("1. bookabledate - 验证日期必须在今天之后")
	fmt.Println("2. chinesephone - 验证中国大陆手机号格式（11位，以1开头）")
	fmt.Println("3. strongpassword - 验证密码强度（至少8位，包含大小写字母和数字）")
	fmt.Println("")
	fmt.Println("测试命令：")
	fmt.Println("")
	fmt.Println("1. 测试预订日期验证（有效日期）：")
	fmt.Println("   curl \"http://localhost:8080/bookings?check_in=2025-12-01&check_out=2025-12-05\"")
	fmt.Println("")
	fmt.Println("2. 测试预订日期验证（无效日期 - 过去的日期）：")
	fmt.Println("   curl \"http://localhost:8080/bookings?check_in=2020-01-01&check_out=2020-01-05\"")
	fmt.Println("")
	fmt.Println("3. 测试用户注册（有效数据）：")
	fmt.Println("   curl -X POST http://localhost:8080/register \\")
	fmt.Println("     -H \"Content-Type: application/json\" \\")
	fmt.Println("     -d '{\"username\":\"zhangsan\",\"email\":\"zhangsan@example.com\",\"phone\":\"13800138000\",\"age\":25,\"password\":\"Password123\"}'")
	fmt.Println("")
	fmt.Println("4. 测试用户注册（无效手机号）：")
	fmt.Println("   curl -X POST http://localhost:8080/register \\")
	fmt.Println("     -H \"Content-Type: application/json\" \\")
	fmt.Println("     -d '{\"username\":\"zhangsan\",\"email\":\"zhangsan@example.com\",\"phone\":\"12345\",\"age\":25,\"password\":\"Password123\"}'")
	fmt.Println("")
	fmt.Println("5. 测试用户注册（弱密码）：")
	fmt.Println("   curl -X POST http://localhost:8080/register \\")
	fmt.Println("     -H \"Content-Type: application/json\" \\")
	fmt.Println("     -d '{\"username\":\"zhangsan\",\"email\":\"zhangsan@example.com\",\"phone\":\"13800138000\",\"age\":25,\"password\":\"weak\"}'")
	fmt.Println("")
	fmt.Println("========================================")
	fmt.Println("")

	// 启动服务器
	router.Run(":8080")
}
