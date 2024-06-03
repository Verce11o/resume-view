package gorilla

import (
	"net/http"

	gorillaMux "github.com/gorilla/mux"
)

type WrappedRouter struct {
	router *gorillaMux.Router
}

func NewWrappedRouter(router *gorillaMux.Router) *WrappedRouter {
	return &WrappedRouter{router: router}
}

func (r *WrappedRouter) MethodFunc(method, path string, handler http.HandlerFunc) {
	r.router.HandleFunc(path, handler).Methods(method)
}

func (r *WrappedRouter) Use(middleware ...func(http.Handler) http.Handler) {
	for _, mw := range middleware {
		r.router.Use(mw)
	}
}

func (r *WrappedRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
