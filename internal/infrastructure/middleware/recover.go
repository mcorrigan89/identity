package middleware

import (
	"log"
	"net/http"
)

func (m *middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				w.Header().Set("Connect", "close")

				log.Fatalf("Panic not recovered %v \n", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
