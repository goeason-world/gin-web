package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 1. router.Static()
	// 将 URL 路径 /assets 映射到本地文件系统目录 ./assets
	// 例如: 访问 /assets/index.html 会提供 ./assets/index.html 文件
	router.Static("/assets", "./assets")

	// 2. router.StaticFS()
	// 与 router.Static() 类似，但可以使用任何实现了 http.FileSystem 接口的源
	// 这里我们使用 http.Dir() 来创建一个文件系统
	router.StaticFS("/more_static", http.Dir("my_file_system"))

	// 3. router.StaticFile()
	// 将一个特定的 URL 路径映射到一个特定的本地文件
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 打印测试说明
	println("=========================================")
	println("静态文件服务示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 router.Static()：")
	println("   curl http://localhost:8080/assets/index.html")
	println("")
	println("2. 测试 router.StaticFile()：")
	println("   curl http://localhost:8080/favicon.ico")
	println("")
	println("3. 测试 router.StaticFS()（需要先在 my_file_system 目录中放入文件）：")
	println("   echo 'Hello from StaticFS' > my_file_system/hello.txt")
	println("   curl http://localhost:8080/more_static/hello.txt")
	println("")
	println("=========================================")
	println("")

	// 启动服务器
	router.Run(":8080")
}
