package gos

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	httpRouter *httprouter.Router
	handlers   []HandlerFunc
	gos        *Gos
	prefix     string
}

func NewRouter(g *Gos) *Router {
	return &Router{
		httpRouter: httprouter.New(),
		gos:        g,
		prefix:     "/",
	}
}

func (r *Router) Use(middlewares ...HandlerFunc) {
	for _, middleware := range middlewares {
		r.handlers = append(r.handlers, middleware)
	}

}
func (r *Router) Group(path string, handlers ...HandlerFunc) *Router {
	handlers = r.combineHandles(handlers)
	return &Router{
		handlers: handlers,
		prefix:   r.path(path),
		gos:      r.gos,
	}
}

//GET handle GET method
func (r *Router) GET(path string, handlers ...HandlerFunc) {
	r.Handle("GET", path, handlers)
}

//POST handle POST method
func (r *Router) POST(path string, handlers ...HandlerFunc) {
	r.Handle("POST", path, handlers)
}

//PATCH handle PATCH method
func (r *Router) PATCH(path string, handlers ...HandlerFunc) {
	r.Handle("PATCH", path, handlers)
}

//PUT handle PUT method
func (r *Router) PUT(path string, handlers ...HandlerFunc) {
	r.Handle("PUT", path, handlers)
}

//DELETE handle DELETE method
func (r *Router) DELETE(path string, handlers ...HandlerFunc) {
	r.Handle("DELETE", path, handlers)
}

//HEAD handle HEAD method
func (r *Router) HEAD(path string, handlers ...HandlerFunc) {
	r.Handle("HEAD", path, handlers)
}

func (r *Router) Handle(method, path string, handlers []HandlerFunc) {
	path = r.path(path)
	handlers = r.combineHandles(handlers)
	r.httpRouter.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		c := r.gos.pool.Get().(*Context)

		fmt.Println(c, "xxx----")

		c.Reset(w, req)
		c.handles = r.handlers
		c.Next()
		r.gos.pool.Put(c)

	})
}

func (r *Router) path(p string) string {
	if r.prefix == "/" {
		return p
	}
	return concat(r.prefix, p)
}

func (r *Router) combineHandles(handlers []HandlerFunc) []HandlerFunc {
	aLen := len(r.handlers)
	hLen := len(handlers)
	h := make([]HandlerFunc, aLen+hLen)
	copy(h, r.handlers)
	for i := 0; i < hLen; i++ {
		h[aLen+i] = handlers[i]
	}
	return h
}

func concat(s ...string) string {
	size := 0
	for i := 0; i < len(s); i++ {
		size += len(s[i])
	}
	buf := make([]byte, 0, size)
	for i := 0; i < len(s); i++ {
		buf = append(buf, []byte(s[i])...)
	}

	return string(buf)
}
