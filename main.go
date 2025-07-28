package main

import (
	"log"
	"net/http"
	"time"
	"web"
)

// group1中间件,令其错误
func OnlyForG1() web.HandleFunc {
	return func(context *web.Context) {
		t := time.Now()
		context.Fail(500, "Internet Server Error")
		log.Printf("[%d] %s in %v for group g1", context.StatusCode, context.Req.RequestURI, time.Since(t))
	}
}
func main() {
	//创造web实例
	webserver := web.New()
	//添加日志中间件
	webserver.Use(web.Logger())
	//添加路由
	webserver.GET("/", func(context *web.Context) {
		context.HTML(http.StatusOK, "<h1>Hello Web</h1>")
	})

	webserver.GET("/hello", func(context *web.Context) {
		context.String(http.StatusOK, "Hello %s, you're at %s\n", context.Query("name"), context.Path)
	})
	webserver.GET("/hello/:name", func(context *web.Context) {
		context.String(http.StatusOK, "Hello %s,you're at %s\n", context.Param("name"), context.Path)
	})
	webserver.GET("/assets/*filepath", func(context *web.Context) {
		context.JSON(http.StatusOK, web.H{
			"filepath": context.Param("filepath"),
		})
	})
	webserver.POST("/login", func(context *web.Context) {
		context.JSON(http.StatusOK, web.H{
			"username": context.PostForm("username"),
			"password": context.PostForm("password"),
		})
	})
	//分组控制
	g1 := webserver.Group("/index")
	g1.Use(OnlyForG1())
	g1.GET("/user/:name", func(context *web.Context) {
		context.String(http.StatusOK, "Hello %s,you're at %s\n", context.Param("name"), context.Path)
	})
	//启动WEB服务
	webserver.Run(":9999")
}
