package core

import "net/http"

type HTTPMixin interface {
	http.Handler
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}
