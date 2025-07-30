package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// 上下文
type Context struct {
	//origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	//request info
	Path   string
	Method string
	Params map[string]string
	//response info
	StatusCode int
	//middlewares
	Middlewares []HandleFunc
	index       int
	//engine
	engine *Engine
}

// Context构造函数
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

// 执行中间件
func (context *Context) Next() {
	context.index++
	for ; context.index < len(context.Middlewares); context.index++ {
		context.Middlewares[context.index](context)
	}
}

// 错误
func (context *Context) Fail(code int, err string) {
	context.index = len(context.Middlewares)
	context.JSON(code, H{"message": err})
}

// 解析动态路由中key匹配值参数
func (context *Context) Param(key string) string {
	value, _ := context.Params[key]
	return value
}

// 提供访问PostForm参数方法
func (context *Context) PostForm(key string) string {
	return context.Req.FormValue(key)
}

// 提供访问Query参数方法
func (context *Context) Query(key string) string {
	return context.Req.URL.Query().Get(key)
}

// 设置状态码
func (context *Context) Status(code int) {
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

// 设置响应头
func (context *Context) SetHeader(key string, value string) {
	context.Writer.Header().Set(key, value)
}

// 快速构造String响应
func (context *Context) String(code int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.Status(code)
	context.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 快速构造JSON响应
func (context *Context) JSON(code int, obj interface{}) {
	context.SetHeader("Content-Type", "application/json")
	context.Status(code)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(context.Writer, err.Error(), 500)
	}
}

// 快速构造Data响应
func (context *Context) Data(code int, data []byte) {
	context.Status(code)
	context.Writer.Write(data)
}

// 快速构造HTML响应
func (context *Context) HTML(code int, name string, data interface{}) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	if err := context.engine.htmlTemplates.ExecuteTemplate(context.Writer, name, data); err != nil {
		context.Fail(500, err.Error())
	}
}
