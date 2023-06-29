package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

var KeyVal = make(map[string]string)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	value, ok := KeyVal[key]
	if !ok {
		w.Write([]byte("Key not found"))
		return
	}

	w.Write([]byte(value))
}

func Set(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}

	bodyString := string(body)

	KeyVal[key] = bodyString
}

func Delete(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	delete(KeyVal, key)
}
