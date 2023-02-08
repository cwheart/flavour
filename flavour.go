package flavour

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type (
	Router struct {
		items map[string]HandleFunc
	}

	Flavour struct {
		server      *http.Server
		router      *Router
		contextPool sync.Pool
	}

	Handler struct {
		f *Flavour
	}

	Context struct {
		w http.ResponseWriter
		r *http.Request
	}

	HandleFunc func(*Context) error
)

func (r *Router) Match(m, uri string) HandleFunc {
	return r.items[fmt.Sprintf("%s:%s", m, uri)]
}

func (r *Router) Add(name string, h HandleFunc) {
	r.items[name] = h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.f.NewContext(w, r)
	defer h.f.contextPool.Put(ctx)
	handleFunc := h.f.Router().Match(r.Method, r.RequestURI)
	if handleFunc == nil {
		w.WriteHeader(404)
		w.Write([]byte("not found"))
	} else {
		handleFunc(ctx)
	}
}

func (f *Flavour) NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := f.contextPool.Get().(*Context)
	ctx.w = w
	ctx.r = r
	return ctx
}

func (c *Context) JSON(code int, body any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(code)
	_, err = c.w.Write(b)
	return err
}

func (c *Context) Param(key string) any {
	return nil
}

func (f *Flavour) Start() error {
	return f.server.ListenAndServe()
}

func (f *Flavour) Router() *Router {
	if f.router == nil {
		f.router = &Router{
			items: make(map[string]HandleFunc),
		}
	}
	return f.router
}

func (f *Flavour) Get(path string, h HandleFunc) {
	f.Router().Add(fmt.Sprintf("GET:%s", path), h)
}

func (f *Flavour) Post(path string, h HandleFunc) {
	f.Router().Add(fmt.Sprintf("POST:%s", path), h)
}

func (f *Flavour) Shutdown(c context.Context) error {
	return f.server.Shutdown(c)
}

func New() *Flavour {
	f := &Flavour{
		contextPool: sync.Pool{
			New: func() any {
				return new(Context)
			},
		},
	}
	f.server = &http.Server{
		Addr: ":8080",
		Handler: &Handler{
			f: f,
		},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return f
}
