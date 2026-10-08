package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello from Snippetbox"))
}

func main() {
	fmt.Println("hello")

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Print("Starting web server")
	err := http.ListenAndServe(":8080", mux)

	log.Fatal(err)

}
