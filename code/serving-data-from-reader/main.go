package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 创建一个默认的 Gin 引擎
	router := gin.Default()

	// 2. 定义一个 GET 路由
	router.GET("/someDataFromReader", func(c *gin.Context) {
		// 3. 从一个 URL 获取数据。这里我们获取的是 Gin 的 Logo 图片。
		// 在实际应用中，数据源可以是任何实现了 io.Reader 的对象，
		// 例如一个本地文件 (os.File)、一个网络连接、或一个正在进行的计算。
		response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		if err != nil || response.StatusCode != http.StatusOK {
			// 如果获取失败，返回 503 服务不可用状态
			c.Status(http.StatusServiceUnavailable)
			return
		}

		// 4. 从 HTTP 响应中提取所需的信息
		// response.Body 是一个 io.ReadCloser，它实现了 io.Reader 接口
		reader := response.Body
		// 获取响应内容的长度
		contentLength := response.ContentLength
		// 获取响应内容的 MIME 类型
		contentType := response.Header.Get("Content-Type")

		// 5. 设置额外的响应头
		// "Content-Disposition" 头告诉浏览器如何处理响应
		// "attachment" 表示将其作为附件处理，触发文件下载
		// "filename" 指定下载后的文件名
		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		// 6. 使用 c.DataFromReader() 将数据流式传输给客户端
		// Gin 会从 reader 中读取数据并将其写入响应体，直到读取完毕。
		// 这种方式避免了将整个文件加载到内存中，对于大文件非常高效。
		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})

	// 7. 启动 HTTP 服务器并监听在 8080 端口
	router.Run(":8080")
}