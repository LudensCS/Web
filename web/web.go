package web

import (
	"net/http"
	"strings"
)

// 定义路由映射的处理方法
type HandleFunc func(context *Context)

// Engine实现ServeHTTP的接口,添加路由映射表
type Engine struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup
}

// Engine构造函数
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
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
	middlewares := make([]HandleFunc, 0)
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.Middlewares...)
		}
	}
	context := NewContext(w, r)
	context.Middlewares = middlewares
	engine.router.handle(context)
}

// 路由分组
type RouterGroup struct {
	prefix      string       //组内前缀
	Middlewares []HandleFunc //中间件
	parent      *RouterGroup //父节点
	engine      *Engine      //共用Engine
}

// 新建分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 按组添加路由
func (group *RouterGroup) AddRoute(method string, comp string, handler HandleFunc) {
	pattern := group.prefix + comp
	group.engine.router.AddRoute(method, pattern, handler)
}

// 按组添加GET路由
func (group *RouterGroup) GET(comp string, handler HandleFunc) {
	group.AddRoute("GET", comp, handler)
}

// 按组添加POST路由
func (group *RouterGroup) POST(comp string, handler HandleFunc) {
	group.AddRoute("POST", comp, handler)
}

// 添加中间件
func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.Middlewares = append(group.Middlewares, middlewares...)
}
