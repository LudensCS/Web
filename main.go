package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
	"web"
)

type student struct {
	Name string
	Age  int
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
func main() {
	//创造web实例
	webserver := web.New()
	//添加日志中间件
	webserver.Use(web.Logger())
	webserver.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	webserver.LoadHTMLGlob("templates/*")
	webserver.Static("/assets", "./static")
	st1 := &student{"zhangsan", 22}
	st2 := &student{"lisi", 25}
	webserver.GET("/", func(context *web.Context) {
		context.HTML(http.StatusOK, "css.tmpl", nil)
	})
	webserver.GET("/students", func(context *web.Context) {
		context.HTML(http.StatusOK, "arr.tmpl", web.H{
			"title":       "web",
			"studentList": [2]*student{st1, st2},
		})
	})
	webserver.GET("/date", func(context *web.Context) {
		context.HTML(http.StatusOK, "custom_func.tmpl", web.H{
			"title": "web",
			"now":   time.Date(2025, 7, 30, 21, 5, 0, 0, time.UTC),
		})
	})
	//启动WEB服务
	webserver.Run(":9999")
}
