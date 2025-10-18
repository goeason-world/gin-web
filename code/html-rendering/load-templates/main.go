package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 使用 LoadHTMLGlob 加载 templates 目录下的所有 .tmpl 文件
	// "templates/**/*.tmpl" 是一个 glob 模式：
	// - templates/ 是基础目录
	// - **/ 表示匹配所有子目录（包括嵌套的子目录）
	// - *.tmpl 表示匹配所有 .tmpl 文件
	// 这样可以一次性加载所有模板文件，包括子目录中的文件
	router.LoadHTMLGlob("templates/**/*.tmpl")

	// 也可以使用 LoadHTMLFiles() 明确指定要加载的文件：
	// router.LoadHTMLFiles("templates/index.tmpl", "templates/posts/index.tmpl", "templates/users/index.tmpl")

	// 定义路由：渲染根目录下的 index.tmpl
	router.GET("/index", func(c *gin.Context) {
		// c.HTML() 用于渲染 HTML 模板
		// 第一个参数：HTTP 状态码
		// 第二个参数：模板文件名（相对于 templates 目录）
		// 第三个参数：传递给模板的数据（gin.H 是 map[string]interface{} 的快捷方式）
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	// 定义路由：渲染 posts 子目录下的 index.tmpl
	router.GET("/posts/index", func(c *gin.Context) {
		// 注意：模板名称需要包含相对于 templates 目录的完整路径
		// 这样可以区分不同目录下的同名模板
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})

	// 定义路由：渲染 users 子目录下的 index.tmpl
	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})

	// 打印测试说明
	println("========================================")
	println("HTML 渲染 - 加载模板示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 使用 LoadHTMLGlob 批量加载模板")
	println("- 支持子目录中的模板文件")
	println("- 处理不同目录下的同名模板")
	println("")
	println("测试命令：")
	println("")
	println("1. 访问主页模板：")
	println("   curl http://localhost:8080/index")
	println("")
	println("2. 访问 posts 模板：")
	println("   curl http://localhost:8080/posts/index")
	println("")
	println("3. 访问 users 模板：")
	println("   curl http://localhost:8080/users/index")
	println("")
	println("4. 或在浏览器中打开：")
	println("   http://localhost:8080/index")
	println("   http://localhost:8080/posts/index")
	println("   http://localhost:8080/users/index")
	println("")
	println("========================================")
	println("")

	// 监听 8080 端口
	router.Run(":8080")
}
