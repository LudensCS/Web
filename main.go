package main

import (
	"fmt"
	"net/http"
	"web"
)

func main() {
	//创造web实例
	webserver := web.New()
	//添加路由
	webserver.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})

	webserver.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	//启动WEB服务
	webserver.Run(":9999")
}
