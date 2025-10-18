package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 定义一个 POST 路由，同时处理查询字符串和表单参数
	router.POST("/post", func(c *gin.Context) {
		// 从查询字符串获取参数
		// 例如: POST /post?id=123&page=1
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")

		// 从 POST 表单获取参数
		// 例如: name=alice&message=hello
		name := c.PostForm("name")
		message := c.PostForm("message")

		// 打印接收到的参数
		fmt.Printf("Query - id: %s, page: %s\n", id, page)
		fmt.Printf("PostForm - name: %s, message: %s\n", name, message)

		// 返回响应
		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"id":      id,
			"page":    page,
			"name":    name,
			"message": message,
		})
	})

	// 使用 GetQuery 和 GetPostForm 检查参数是否存在
	router.POST("/check", func(c *gin.Context) {
		// GetQuery 返回两个值：参数值和是否存在
		id, idExists := c.GetQuery("id")
		page, pageExists := c.GetQuery("page")

		// GetPostForm 返回两个值：参数值和是否存在
		name, nameExists := c.GetPostForm("name")
		message, messageExists := c.GetPostForm("message")

		// 构建响应
		response := gin.H{
			"query": gin.H{
				"id": gin.H{
					"value":  id,
					"exists": idExists,
				},
				"page": gin.H{
					"value":  page,
					"exists": pageExists,
				},
			},
			"form": gin.H{
				"name": gin.H{
					"value":  name,
					"exists": nameExists,
				},
				"message": gin.H{
					"value":  message,
					"exists": messageExists,
				},
			},
		}

		c.JSON(http.StatusOK, response)
	})

	// 实际应用示例：搜索功能
	router.POST("/search", func(c *gin.Context) {
		// 从查询字符串获取分页参数
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("page_size", "10")

		// 从 POST 表单获取搜索条件
		keyword := c.PostForm("keyword")
		category := c.DefaultPostForm("category", "all")
		sortBy := c.DefaultPostForm("sort", "relevance")

		fmt.Printf("Search - keyword: %s, category: %s, sort: %s, page: %s, page_size: %s\n",
			keyword, category, sortBy, page, pageSize)

		c.JSON(http.StatusOK, gin.H{
			"message":   "搜索成功",
			"keyword":   keyword,
			"category":  category,
			"sort":      sortBy,
			"page":      page,
			"page_size": pageSize,
			"results":   []string{"结果1", "结果2", "结果3"},
		})
	})

	// 实际应用示例：用户注册
	router.POST("/register", func(c *gin.Context) {
		// 从查询字符串获取推荐人 ID（可选）
		referrer := c.DefaultQuery("ref", "")

		// 从 POST 表单获取用户信息
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		// 验证必需字段
		if username == "" || email == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "用户名、邮箱和密码不能为空",
			})
			return
		}

		fmt.Printf("Register - username: %s, email: %s, referrer: %s\n",
			username, email, referrer)

		c.JSON(http.StatusOK, gin.H{
			"message":  "注册成功",
			"username": username,
			"email":    email,
			"referrer": referrer,
		})
	})

	// 打印测试说明
	println("=========================================")
	println("查询字符串和表单参数示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 使用 c.Query() 获取查询字符串参数")
	println("- 使用 c.PostForm() 获取 POST 表单参数")
	println("- 使用 c.DefaultQuery() 和 c.DefaultPostForm() 设置默认值")
	println("- 使用 c.GetQuery() 和 c.GetPostForm() 检查参数是否存在")
	println("")
	println("测试命令：")
	println("")
	println("1. 基本示例（同时使用查询字符串和表单参数）：")
	println("   curl -X POST 'http://localhost:8080/post?id=123&page=1' \\")
	println("     -d 'name=alice&message=hello'")
	println("")
	println("2. 检查参数是否存在：")
	println("   curl -X POST 'http://localhost:8080/check?id=123' \\")
	println("     -d 'name=bob'")
	println("")
	println("3. 搜索功能示例：")
	println("   curl -X POST 'http://localhost:8080/search?page=2&page_size=20' \\")
	println("     -d 'keyword=golang&category=tutorial&sort=date'")
	println("")
	println("4. 用户注册示例：")
	println("   curl -X POST 'http://localhost:8080/register?ref=user123' \\")
	println("     -d 'username=newuser&email=user@example.com&password=secret'")
	println("")
	println("预期输出：")
	println("- 每个请求都会返回包含查询字符串和表单参数的 JSON 响应")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
