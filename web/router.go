package web

import (
	"log"
	"strings"
)

type Router struct {
	roots    map[string]*node
	Handlers map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		Handlers: make(map[string]HandleFunc),
	}
}

// 解析路径为[]string
func ParsePattern(pattern string) []string {
	sPattern := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, part := range sPattern {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
func (router *Router) AddRoute(method string, pattern string, handler HandleFunc) {
	log.Printf("Route %6s - %s", method, pattern)
	parts := ParsePattern(pattern)
	_, ok := router.roots[method]
	if !ok {
		router.roots[method] = &node{}
	}
	router.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern
	router.Handlers[key] = handler
}

// 调用处理函数
func (router *Router) handle(context *Context) {
	nd, params := router.getRoute(context.Method, context.Path)
	if nd != nil {
		context.Params = params
		key := context.Method + "-" + nd.pattern
		router.Handlers[key](context)
	} else {
		context.String(404, "404 NOT FOUND: %s\n", context.Path)
	}
}

// 匹配路由并解析匹配参数
func (router *Router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := ParsePattern(pattern)
	params := make(map[string]string)
	root, ok := router.roots[method]
	if !ok {
		return nil, nil
	}
	nd := root.search(searchParts, 0)
	if nd == nil {
		return nil, nil
	}
	parts := ParsePattern(nd.pattern)
	for index, part := range parts {
		if part[0] == ':' && len(part) > 1 {
			params[part[1:]] = searchParts[index]
		} else if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}
	return nd, params
}
