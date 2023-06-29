package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	serverPort  string
	StoragePath = "./keyVal"
	key_val     = make(map[string]string) // for now, we will handle complex values
)

const url = "http://localhost"

func main() {
	// defer fmt.Println("Program terminated")
	// defer saveFile(key_val)
	serverPort = "8080" // if env port is not set

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		serverPort = fromEnv
	}

	log.Printf("Satrting up on %s:%s", url, serverPort)
	// Set routing rule

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// read from the file, if it exsists and add it into th key-val
	// if it does not exist, create it using the saveFile func

	file, err := os.OpenFile(StoragePath+"/keyVal.json", os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// decode the json into the map
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&key_val)
	if err != nil {
		log.Fatalf("Failed to decode data: %v", err)
	}

	r.Get("/", home)
	r.Get("/key/{key}", Get)
	r.Post("/key/{key}", Set)
	r.Delete("/key/{key}", Delete)

	// get port from env

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
		saveFile(file, key_val)
		os.Exit(0)
	}()

	for {
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

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

func saveFile(f *os.File, key_val map[string]string) {
	//write the json into the file
	encodedData, err := json.Marshal(key_val)
	if err != nil {
		log.Fatalf("Failed to encode data: %v", err)
	}

	f.Write(encodedData)

}
