package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadForm 定义了包含文件上传的表单结构
// 使用 form 标签来指定表单字段名
// 使用 binding:"required" 标签来指定字段为必需
type UploadForm struct {
	User        string                `form:"user" binding:"required"`
	Description string                `form:"description"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 为 /upload 路由设置一个较低的内存限制（默认是 32 MiB）
	// 这可以防止恶意用户通过上传大文件来耗尽服务器内存
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 定义一个 POST 路由 /upload
	// 这个端点用于处理文件上传
	router.POST("/upload", func(c *gin.Context) {
		var form UploadForm

		// c.ShouldBind() 会自动处理 multipart/form-data
		// 它会绑定表单字段和上传的文件
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 确保 uploads 目录存在
		// 在实际应用中，应该在程序启动时检查或创建
		// os.MkdirAll("./uploads", 0755)

		// 保存上传的文件到服务器上的指定位置
		// 这里我们保存到 ./uploads/ 目录下
		// c.SaveUploadedFile() 是一个方便的方法，可以处理文件保存
		filePath := "./uploads/" + form.File.Filename
		if err := c.SaveUploadedFile(form.File, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
			return
		}

		// 打印日志，方便调试
		log.Printf("文件上传成功: user=%s, filename=%s, size=%d", form.User, form.File.Filename, form.File.Size)

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"message":     "上传成功",
			"user":        form.User,
			"description": form.Description,
			"filename":    form.File.Filename,
			"size":        form.File.Size,
			"path":        filePath,
		})
	})

	// 打印测试说明
	fmt.Println("=========================================")
	fmt.Println("文件上传扩展示例")
	fmt.Println("=========================================")
	fmt.Println("服务器启动在 :8080")
	fmt.Println("")
	fmt.Println("功能说明：")
	fmt.Println("- 演示如何处理包含文件上传的 multipart/form-data 表单")
	fmt.Println("- 将表单字段和文件绑定到结构体")
	fmt.Println("- 保存上传的文件到服务器")
	fmt.Println("")
	fmt.Println("测试命令：")
	fmt.Println("")
	fmt.Println("1. 创建一个测试文件：")
	fmt.Println("   echo \"hello world\" > test.txt")
	fmt.Println("")
	fmt.Println("2. 创建 uploads 目录：")
	fmt.Println("   mkdir -p uploads")
	fmt.Println("")
	fmt.Println("3. 发送文件上传请求：")
	fmt.Println("   curl -X POST http://localhost:8080/upload \\")
	fmt.Println("     -F \"user=testuser\" \\")
	fmt.Println("     -F \"description=这是一个测试文件\" \\")
	fmt.Println("     -F \"file=@test.txt\"")
	fmt.Println("")
	fmt.Println("预期输出：")
	fmt.Println("- 客户端将收到包含文件信息的 JSON 响应")
	fmt.Println("- 服务器将在 uploads 目录下保存 test.txt 文件")
	fmt.Println("=========================================")
	fmt.Println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
