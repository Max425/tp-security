package proxy

import (
	"log"
	"net/http"
	"runtime/debug"
)

type Middleware struct {
	next http.Handler
}

func NewMiddleware(next http.Handler) *Middleware {
	return &Middleware{next: next}
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error: %s %s", err.(string), string(debug.Stack()))
		}
	}()
	log.Printf("Received request: %s %s %s", r.Method, r.URL.Path, r.Proto)

	m.next.ServeHTTP(w, r)
}
