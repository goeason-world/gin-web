package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 构建带有 Logger 和 Recovery 中间件的默认引擎，便于观察和排错。
	router := gin.Default()

	// 2. 限制 multipart 表单使用的内存，避免大文件长时间堆积在内存中。
	// Gin 在超过限制后会将剩余部分写入临时文件，从而保护服务。
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 3. 确保文件保存目录存在；MkdirAll 会在目录已存在时静默返回。
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Fatalf("create upload dir %q failed: %v", uploadDir, err)
	}

	// 4. 注册处理单文件上传的路由。
	router.POST("/upload", func(c *gin.Context) {
		// 5. 从表单字段中抓取文件对象；字段名需要与 HTML 表单的 name 属性一致。
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err))
			return
		}

		// 6. 记录文件名和大小，方便排错或审计。
		log.Printf("received file: %s (%d bytes)", file.Filename, file.Size)

		// 7. 组装目标路径并保存文件；SaveUploadedFile 内部会自动打开源文件并写入目标文件。
		dst := fmt.Sprintf("%s/%s", uploadDir, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err))
			return
		}

		// 8. 返回明确的上传成功响应。
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// 9. 启动 HTTP 服务；Run 在内部调用 http.ListenAndServe，并会返回任何运行时错误。
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
