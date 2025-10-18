package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// secrets 模拟业务数据；键为用户名，值为需要保护的信息。
var secrets = map[string]gin.H{
	"alice": {
		"email": "alice@example.com",
		"phone": "123-456",
	},
	"bob": {
		"email": "bob@example.com",
		"phone": "987-654",
	},
}

func main() {
	// gin.Default() 默认挂载 Logger 和 Recovery，中间件会记录访问日志并在 panic 时恢复。
	router := gin.Default()

	// 将需要认证的路由划分到 /admin 分组；gin.BasicAuth 负责校验 Authorization 头。
	// gin.Accounts 是一个用户名到密码的映射，这里定义了两个演示账号。
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"alice": "wonderland",
		"bob":   "builder",
	}))

	// GET /admin/secrets 仅当 BasicAuth 校验通过后才会执行 handler。
	authorized.GET("/secrets", func(c *gin.Context) {
		// gin.AuthUserKey 存放认证通过的用户名，可用于精细化权限控制。
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			// 返回当前用户对应的敏感数据；如果需要可以继续拼接业务字段。
			c.JSON(http.StatusOK, gin.H{
				"user":   user,
				"secret": secret,
			})
			return
		}

		// 用户认证通过但没有配置数据时，明确返回 404，避免暴露其他用户信息。
		c.JSON(http.StatusNotFound, gin.H{
			"message": "no secret for this user",
		})
	})

	// 监听 8080 端口；生产环境可通过环境变量或配置文件覆盖。
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
