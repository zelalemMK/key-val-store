package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var serverPort string

var key_val = make(map[string]string) // for now, we will handle complex values

const url = "http://localhost"

func main() {
	serverPort = "8080" // if env port is not set

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		serverPort = fromEnv
	}

	log.Printf("Satrting up on %s:%s", url, serverPort)
	// Set routing rule

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", home)
	r.Get("/key/{key}", Get)
	r.Post("/key/{key}", Set)
	r.Delete("/key/{key}", Delete)

	// get port from env

	err := http.ListenAndServe(":"+serverPort, r)
	if err != nil {
		log.Fatal("Error running the server")
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// r.Get("/key/{key}")
func Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	value, ok := key_val[key]
	if !ok {
		w.Write([]byte("Key not found"))
		return
	}

	w.Write([]byte(value))
}

func Set(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	fmt.Println(key)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}

	bodyString := string(body)

	key_val[key] = bodyString
	// w.Write([]byte(key))
}

func Delete(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	delete(key_val, key)
}
