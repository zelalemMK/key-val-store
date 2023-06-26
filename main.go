package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var serverPort string

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

	// get port from env

	err := http.ListenAndServe(":"+serverPort, r)
	if err != nil {
		log.Fatal("Error running the server")
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
