package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化默认引擎，包含访问日志与 panic 自动恢复。
	router := gin.Default()

	// 2. 限制 multipart 表单在内存中的占用，过大的请求 Gin 会落盘处理。
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 3. 创建上传目录；使用 MkdirAll 以便多级目录也能一次性创建。
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Fatalf("create upload dir %q failed: %v", uploadDir, err)
	}

	// 4. 处理多文件上传；表单字段需设置为 <input type="file" name="upload[]" multiple>
	router.POST("/upload", func(c *gin.Context) {
		// 5. 解析整个 multipart 表单，便于一次性获取所有文件。
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err))
			return
		}

		// 6. files 包含表单字段 upload[] 上传的所有文件。
		files := form.File["upload[]"]
		if len(files) == 0 {
			c.String(http.StatusBadRequest, "no file received")
			return
		}

		// 7. 依次保存每一个文件，过程中如遇错误立即返回以保持客户端知情。
		for _, file := range files {
			log.Printf("received file: %s (%d bytes)", file.Filename, file.Size)

			dst := fmt.Sprintf("%s/%s", uploadDir, file.Filename)
			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err))
				return
			}
		}

		// 8. 汇报成功结果，包含已处理的文件数量。
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	// 9. 以 8080 端口启动 HTTP 服务。
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
