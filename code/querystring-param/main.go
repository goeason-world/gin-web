package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的 Gin 引擎
	router := gin.Default()

	// 基本示例：使用 Query 和 DefaultQuery
	// 访问: /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		// 使用 c.Query() 获取参数，不存在返回空字符串
		firstname := c.Query("firstname")

		// 使用 c.DefaultQuery() 获取参数，不存在返回默认值
		lastname := c.DefaultQuery("lastname", "Guest")

		c.JSON(http.StatusOK, gin.H{
			"message":   "Welcome",
			"firstname": firstname,
			"lastname":  lastname,
		})
	})

	// 使用 GetQuery 检查参数是否存在
	// 访问: /check?name=alice
	router.GET("/check", func(c *gin.Context) {
		// GetQuery 返回两个值：参数值和是否存在
		name, exists := c.GetQuery("name")

		if exists {
			c.JSON(http.StatusOK, gin.H{
				"message": "参数存在",
				"name":    name,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "参数不存在",
			})
		}
	})

	// 处理数组参数
	// 访问: /array?ids=1&ids=2&ids=3
	router.GET("/array", func(c *gin.Context) {
		// QueryArray 获取同名参数的数组
		ids := c.QueryArray("ids")

		fmt.Printf("IDs: %v\n", ids)

		c.JSON(http.StatusOK, gin.H{
			"message": "接收到数组参数",
			"ids":     ids,
			"count":   len(ids),
		})
	})

	// 处理 Map 参数
	// 访问: /map?user[name]=alice&user[age]=20
	router.GET("/map", func(c *gin.Context) {
		// QueryMap 获取 map 形式的参数
		user := c.QueryMap("user")

		fmt.Printf("User: %v\n", user)

		c.JSON(http.StatusOK, gin.H{
			"message": "接收到 Map 参数",
			"user":    user,
		})
	})

	// 实际应用示例：搜索功能
	// 访问: /search?keyword=golang&page=1&page_size=20&sort=date
	router.GET("/search", func(c *gin.Context) {
		keyword := c.Query("keyword")
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("page_size", "10")
		sortBy := c.DefaultQuery("sort", "relevance")

		fmt.Printf("Search - keyword: %s, page: %s, page_size: %s, sort: %s\n",
			keyword, page, pageSize, sortBy)

		c.JSON(http.StatusOK, gin.H{
			"message":   "搜索成功",
			"keyword":   keyword,
			"page":      page,
			"page_size": pageSize,
			"sort":      sortBy,
			"results":   []string{"结果1", "结果2", "结果3"},
		})
	})

	// 实际应用示例：过滤器
	// 访问: /products?category=electronics&min_price=100&max_price=1000&brands=apple&brands=samsung
	router.GET("/products", func(c *gin.Context) {
		category := c.DefaultQuery("category", "all")
		minPrice := c.DefaultQuery("min_price", "0")
		maxPrice := c.DefaultQuery("max_price", "999999")
		brands := c.QueryArray("brands")

		fmt.Printf("Filter - category: %s, price: %s-%s, brands: %v\n",
			category, minPrice, maxPrice, brands)

		c.JSON(http.StatusOK, gin.H{
			"message":   "获取产品列表",
			"category":  category,
			"min_price": minPrice,
			"max_price": maxPrice,
			"brands":    brands,
			"products":  []string{"产品1", "产品2", "产品3"},
		})
	})

	// 打印测试说明
	println("=========================================")
	println("查询字符串参数示例")
	println("=========================================")
	println("服务器启动在 :8080")
	println("")
	println("功能说明：")
	println("- 使用 c.Query() 获取参数")
	println("- 使用 c.DefaultQuery() 设置默认值")
	println("- 使用 c.GetQuery() 检查参数是否存在")
	println("- 使用 c.QueryArray() 获取数组参数")
	println("- 使用 c.QueryMap() 获取 Map 参数")
	println("")
	println("测试命令：")
	println("")
	println("1. 基本示例：")
	println("   curl 'http://localhost:8080/welcome?firstname=Jane&lastname=Doe'")
	println("")
	println("2. 使用默认值：")
	println("   curl 'http://localhost:8080/welcome?firstname=Jane'")
	println("")
	println("3. 检查参数是否存在：")
	println("   curl 'http://localhost:8080/check?name=alice'")
	println("   curl 'http://localhost:8080/check'")
	println("")
	println("4. 数组参数：")
	println("   curl 'http://localhost:8080/array?ids=1&ids=2&ids=3'")
	println("")
	println("5. Map 参数：")
	println("   curl --globoff 'http://localhost:8080/map?user[name]=alice&user[age]=20'")
	println("")
	println("6. 搜索功能：")
	println("   curl 'http://localhost:8080/search?keyword=golang&page=2&page_size=20&sort=date'")
	println("")
	println("7. 过滤器：")
	println("   curl 'http://localhost:8080/products?category=electronics&min_price=100&max_price=1000&brands=apple&brands=samsung'")
	println("")
	println("预期输出：")
	println("- 每个请求都会返回包含查询参数的 JSON 响应")
	println("=========================================")
	println("")

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
