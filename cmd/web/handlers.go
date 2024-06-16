package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (a *app) heyHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	if err := ts.ExecuteTemplate(w, "base", nil); err != nil {
		a.serverError(w, r, err)
	}
}

func (a *app) snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(fmt.Sprintf("Snippet View for ID: %d", id)))
}

func (a *app) snippetFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Snippet Form"))
}

func (a *app) snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Create Snippet"))
}
