package web

import (
	"net/http"
)

// 定义路由映射的处理方法
type HandleFunc func(context *Context)

// Engine实现ServeHTTP的接口,添加路由映射表
type Engine struct {
	router *Router
}

// Engine构造函数
func New() *Engine {
	return &Engine{router: NewRouter()}
}

// 添加静态路由的方法
func (engine *Engine) AddRoute(method string, pattern string, handler HandleFunc) {
	engine.router.AddRoute(method, pattern, handler)
}

// 添加GET请求的方法
func (engine *Engine) GET(pattern string, handler HandleFunc) {
	engine.AddRoute("GET", pattern, handler)
}

// 添加POST请求的方法
func (engine *Engine) POST(pattern string, handler HandleFunc) {
	engine.AddRoute("POST", pattern, handler)
}

// 启动HTTP服务端的方法
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// 解析请求路径,查找路由映射表。如果查到就执行注册的方法,否则返回404
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := NewContext(w, r)
	engine.router.handle(context)
}
