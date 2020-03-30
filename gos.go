package gos

import (
	"html/template"
	"net/http"
	"sync"
)

type (
	HandlerFunc func(c *Context)
	Gos         struct {
		Router   *Router
		Template *template.Template
		pool     sync.Pool
	}
)

func New() *Gos {
	g := &Gos{}
	g.Router = NewRouter(g)
	g.pool.New = func() interface{} {
		return &Context{
			gos:     g,
			index:   -1,
			Reponse: &Response{},
		}
	}
	return g
}

func Default() *Gos {
	g := New()
	g.Router.Use(Recovery())
	g.Router.Use(Logger())

	return g
}

func (g *Gos) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.Router.httpRouter.ServeHTTP(w, r)
}

func (g *Gos) Run(addr string) {
	if err := http.ListenAndServe(addr, g); err != nil {
		panic(err)
	}

}

func (g *Gos) LoadHtml(pattern string) {
	g.Template = template.Must(template.ParseGlob(pattern))
}
