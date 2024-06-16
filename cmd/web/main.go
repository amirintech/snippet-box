package main

import (
	"fmt"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", heyHandler)
	mux.HandleFunc("GET /snippet/view/{id}", snippetViewHandler)
	mux.HandleFunc("GET /snippet/form", snippetFormHandler)
	mux.HandleFunc("POST /snippet/create", snippetFormHandler)

	addr := ":5001"
	fmt.Printf("Listening on http://localhost%s\n", addr)
	http.ListenAndServe(addr, mux)
}
