package main

import (
	"net/http"
	"web"
)

func main() {
	//创造web实例
	webserver := web.New()
	//添加路由
	webserver.GET("/", func(context *web.Context) {
		context.HTML(http.StatusOK, "<h1>Hello Web</h1>")
	})

	webserver.GET("/hello", func(context *web.Context) {
		context.String(http.StatusOK, "Hello %s, you're at %s\n", context.Query("name"), context.Path)
	})
	webserver.POST("/login", func(context *web.Context) {
		context.JSON(http.StatusOK, web.H{
			"username": context.PostForm("username"),
			"password": context.PostForm("password"),
		})
	})
	//启动WEB服务
	webserver.Run(":9999")
}
