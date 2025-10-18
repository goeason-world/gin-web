package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 日志文件按目录组织，便于在容器或宿主机中挂载。
	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Fatalf("create log dir failed: %v", err)
	}

	// 以追加模式打开日志文件，避免重启服务时覆盖历史日志。
	logFile, err := os.OpenFile("logs/gin-access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("open log file failed: %v", err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Printf("close log file failed: %v", err)
		}
	}()

	// Gin 默认会输出带颜色的日志，写入文件时关闭彩色字符以避免乱码。
	gin.DisableConsoleColor()

	// 同时写入日志文件与标准输出，便于本地开发排查。
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	// 使用 gin.New 以避免自动添加 Logger；我们手动注册以确保输出按照我们配置的 writer。
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 业务路由仅作占位。访问该路由即可触发日志写入。
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
