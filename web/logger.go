package web

import (
	"log"
	"time"
)

// 日志中间件
func Logger() HandleFunc {
	return func(context *Context) {
		t := time.Now()
		context.Next()
		log.Printf("[%d] %s in %v", context.StatusCode, context.Req.RequestURI, time.Since(t))
	}
}
