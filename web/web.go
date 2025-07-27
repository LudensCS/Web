package web

import (
	"net/http"
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
	context := NewContext(w, r)
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

func (group *RouterGroup) AddRoute(method string, comp string, handler HandleFunc) {
	pattern := group.prefix + comp
	group.engine.router.AddRoute(method, pattern, handler)
}
func (group *RouterGroup) GET(comp string, handler HandleFunc) {
	group.AddRoute("GET", comp, handler)
}
func (group *RouterGroup) POST(comp string, handler HandleFunc) {
	group.AddRoute("POST", comp, handler)
}
