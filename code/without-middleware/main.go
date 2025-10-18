package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.New() 返回一个没有任何默认中间件的 Engine；这意味着不会自动记录访问日志，也不会在 panic 时恢复。
	router := gin.New()

	// 如果仍然希望记录请求，需要手动实现。这里演示在 handler 内写入一条简易日志。
	router.GET("/ping", func(c *gin.Context) {
		log.Printf("%s %s -> 200", c.Request.Method, c.Request.URL.Path)
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 没有 Recovery 中间件时，一旦发生 panic 服务会直接崩溃；该路由用于显式观察这一差异。
	router.GET("/panic", func(c *gin.Context) {
		panic("server crashed because Recovery middleware is absent")
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
