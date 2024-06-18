package main

import "net/http"

func (a *app) routes() *http.ServeMux {
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", a.homeHandler)
	mux.HandleFunc("GET /snippet/view/{id}", a.snippetViewHandler)
	mux.HandleFunc("GET /snippet/form", a.snippetFormHandler)
	mux.HandleFunc("POST /snippet/create", a.snippetCreateHandler)

	return mux
}
