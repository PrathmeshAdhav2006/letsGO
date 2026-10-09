package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	infologger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	type config struct {
		addr string
	}

	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errlogger,
		Handler:  mux,
	}

	infologger.Printf("Starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	errlogger.Fatal(err)
}
