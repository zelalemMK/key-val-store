package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/zelalemmk/key-val-store/cmd/web/handlers"
	utils "github.com/zelalemmk/key-val-store/internal/helper"
)

var (
	serverPort string
)

const (
	url        = "http://localhost"
	fileName   = "keyVal.json"
	filePerm   = 0644
	dirPerm    = 0755
	numEntries = 100
)

func init() {
	err := utils.GenerateData(numEntries)
	if err != nil {
		log.Fatalf("Failed to generate data: %v", err)
	}

	log.Println("Sample data generated")

	utils.LoadData()
}

func main() {
	serverPort = "8080" // default port if not set in environment

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		serverPort = fromEnv
	}

	log.Printf("Starting up on %s:%s", url, serverPort)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handlers.Home)
	r.Get("/key/{key}", handlers.Get)
	r.Post("/key/{key}", handlers.Set)
	r.Delete("/key/{key}", handlers.Delete)

	go func() {
		err := http.ListenAndServe(":"+serverPort, r)
		if err != nil {
			log.Fatal("Error running the server")
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("Received an interrupt, stopping services...")
		utils.SaveData()
		os.Exit(0)
	}()

	select {}
}
