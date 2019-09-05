package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1234"
	}

	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, "hello world")
	})

	http.ListenAndServe(":"+port, httplog(secure(app)))
}

func httplog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// secure verifies that all requests have an extremely strong password
func secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// auth the request
		if r.Header.Get("Authorization") != "hunter2" {
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
