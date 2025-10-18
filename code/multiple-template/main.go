package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 设置自定义的模板分隔符
	// 默认的 {{ }} 可能与前端框架（如 Vue.js）冲突
	// 这里改为使用 {[ ]} 作为模板分隔符
	router.Delims("{[", "]}")

	// 设置自定义的模板函数
	// 这些函数可以在模板中直接调用
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})

	// 加载所有模板文件
	// 使用通配符 ** 递归加载 templates 目录下的所有 .tmpl 文件
	router.LoadHTMLGlob("templates/**/*")
	// 或者使用 LoadHTMLFiles 加载指定的文件：
	// router.LoadHTMLFiles("templates/posts/index.tmpl", "templates/users/index.tmpl")

	// 定义 /posts/index 路由
	router.GET("/posts/index", func(c *gin.Context) {
		// 渲染 posts/index.tmpl 模板
		// 第二个参数是传递给模板的数据
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
			"posts": []gin.H{
				{"title": "Go 语言入门", "date": time.Now()},
				{"title": "Gin 框架实战", "date": time.Now().AddDate(0, 0, -1)},
				{"title": "RESTful API 设计", "date": time.Now().AddDate(0, 0, -2)},
			},
		})
	})

	// 定义 /users/index 路由
	router.GET("/users/index", func(c *gin.Context) {
		// 渲染 users/index.tmpl 模板
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
			"users": []gin.H{
				{"name": "张三", "email": "zhangsan@example.com"},
				{"name": "李四", "email": "lisi@example.com"},
				{"name": "王五", "email": "wangwu@example.com"},
			},
		})
	})

	// 打印测试说明
	println("=========================================")
	println("使用多个模板文件示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 从多个目录加载模板文件")
	println("- 使用自定义的模板分隔符 {[ ]}")
	println("- 使用自定义模板函数 formatAsDate")
	println("- 渲染不同的 HTML 页面")
	println("")
	println("测试命令：")
	println("")
	println("1. 访问文章列表页面：")
	println("   curl http://localhost:8080/posts/index")
	println("   或在浏览器中打开：http://localhost:8080/posts/index")
	println("")
	println("2. 访问用户列表页面：")
	println("   curl http://localhost:8080/users/index")
	println("   或在浏览器中打开：http://localhost:8080/users/index")
	println("")
	println("预期输出：")
	println("- 两个页面会显示不同的内容和样式")
	println("- 文章列表会使用自定义函数格式化日期")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
