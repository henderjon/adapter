package main

import (
	"log"
	"net/http"
	"net/textproto"
	"time"
)

func main() {
	allmux := http.NewServeMux()
	allmux.HandleFunc("/", ServeIndex)

	allmuxHandler := Adapt(
		allmux,
		versionAdapter(`greetings`),
	)

	srv := &http.Server{
		Addr:         ":8081",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      allmuxHandler,
	}

	log.Println("listening:", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err, true)
	}
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Hello, World.`))
}

func versionAdapter(greeting string) Adapter { // inject and return adapter
	return Adapter(func(handler http.Handler) http.Handler { // return handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(textproto.CanonicalMIMEHeaderKey(`vnd-adapter-example`), greeting)
			handler.ServeHTTP(w, r)
		})
	})
}
