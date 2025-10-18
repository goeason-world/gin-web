package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StructA 是最基础的结构体，包含一个字段
type StructA struct {
	FieldA string `form:"field_a"`
}

// StructB 包含一个嵌套的 StructA 结构体（值类型）
type StructB struct {
	NestedStruct StructA
	FieldB       string `form:"field_b"`
}

// StructC 包含一个嵌套的 StructA 结构体指针
type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}

// StructD 包含一个匿名嵌套结构体
type StructD struct {
	NestedAnonyStruct struct {
		FieldX string `form:"field_x"`
	}
	FieldD string `form:"field_d"`
}

// GetDataB 处理 StructB 类型的表单数据绑定
func GetDataB(c *gin.Context) {
	var b StructB
	// 绑定表单数据到结构体
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 返回绑定后的数据
	c.JSON(http.StatusOK, gin.H{
		"a": b.NestedStruct,
		"b": b.FieldB,
	})
}

// GetDataC 处理 StructC 类型的表单数据绑定
func GetDataC(c *gin.Context) {
	var b StructC
	// 绑定表单数据到结构体
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 返回绑定后的数据
	c.JSON(http.StatusOK, gin.H{
		"a": b.NestedStructPointer,
		"c": b.FieldC,
	})
}

// GetDataD 处理 StructD 类型的表单数据绑定
func GetDataD(c *gin.Context) {
	var b StructD
	// 绑定表单数据到结构体
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 返回绑定后的数据
	c.JSON(http.StatusOK, gin.H{
		"x": b.NestedAnonyStruct,
		"d": b.FieldD,
	})
}

func main() {
	router := gin.Default()

	// 注册路由
	router.GET("/getb", GetDataB)
	router.GET("/getc", GetDataC)
	router.GET("/getd", GetDataD)

	// 打印测试说明
	println("========================================")
	println("绑定表单数据至自定义结构体示例")
	println("========================================")
	println("服务器启动在 :8080")
	println("")
	println("测试命令：")
	println("")
	println("1. 测试 StructB（嵌套结构体值类型）：")
	println(`   curl "http://localhost:8080/getb?field_a=hello&field_b=world"`)
	println("")
	println("2. 测试 StructC（嵌套结构体指针）：")
	println(`   curl "http://localhost:8080/getc?field_a=hello&field_c=world"`)
	println("")
	println("3. 测试 StructD（匿名嵌套结构体）：")
	println(`   curl "http://localhost:8080/getd?field_x=hello&field_d=world"`)
	println("")
	println("========================================")
	println("")

	// 启动服务器
	router.Run(":8080")
}
