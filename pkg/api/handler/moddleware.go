package handler

import (
	"log"
	"net/http"
	"runtime/debug"
)

func (h *Handler) panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				log.Println("Panic",
					r.Method,
					r.RequestURI,
					err.(string),
					string(debug.Stack()),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
