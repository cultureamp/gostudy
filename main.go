package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RequestThing struct {
	Name string
	Camp string
}

type ResponseThing struct {
	Greeting string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1234"
	}

	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		r.Body.Close()

		var rt RequestThing
		json.Unmarshal(data, &rt)

		respBytes, err := json.Marshal(ResponseThing{Greeting: "hello " + rt.Name})
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(respBytes)
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
