package router

import (
	"net/http"
	"signature/middleware"

	"signature/controller"

	"github.com/gin-gonic/gin"
)

func Router() http.Handler {
	engine := gin.Default()

	// 无需验证签名请求，用于验证通道是否ok
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong by webServer.",
		})
	})

	// 在此之后的接口都需要通过签名验证，才能正常调用
	engine.Use(middleware.SignatureMiddleware())

	// 测试GET请求，无包体签名
	engine.GET("say-hello", controller.SayHello)

	// 测试POST请求，有包体签名
	engine.POST("do-something", controller.DoSomething)

	return engine
}
