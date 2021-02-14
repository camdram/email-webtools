package server

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gorilla/handlers"
)

func logMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func serverHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "https://github.com/camdram/email-webtools")
		next.ServeHTTP(w, r)
	})
}

func (c *Controller) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		expected := "Bearer " + c.bearer
		switch auth {
		case expected:
			next.ServeHTTP(w, r)
		case "":
			w.Header().Set("WWW-Authenticate", "Bearer realm=\"Camdram email-webtools\"")
			http.Error(w, "401 unauthorized", http.StatusUnauthorized)
		default:
			http.Error(w, "403 forbidden", http.StatusForbidden)
		}
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				log.Println("An internal server error occurred:", err)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
