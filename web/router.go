package web

import "log"

type Router struct {
	Handlers map[string]HandleFunc
}

func NewRouter() *Router {
	return &Router{
		Handlers: make(map[string]HandleFunc),
	}
}
func (router *Router) AddRoute(method string, pattern string, handler HandleFunc) {
	log.Printf("Route %6s - %s", method, pattern)
	key := method + "-" + pattern
	router.Handlers[key] = handler
}
func (router *Router) handle(context *Context) {
	key := context.Method + "-" + context.Path
	if handler, ok := router.Handlers[key]; ok {
		handler(context)
	} else {
		context.String(404, "404 NOT FOUND: %s\n", context.Path)
	}
}
