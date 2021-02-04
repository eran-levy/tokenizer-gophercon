package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dummy server"))

}

func main() {
	http.HandleFunc("/language", handler)
	log.Fatal(http.ListenAndServe(":3333", nil))
}
