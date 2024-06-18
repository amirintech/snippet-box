package main

import (
	"errors"
	"fmt"
	"github.com/amirintech/snippet-box/internal/models"
	"net/http"
	"strconv"
)

func (a *app) homeHandler(w http.ResponseWriter, r *http.Request) {
	snippets, err := a.snippetModel.GetLatest()
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	data := a.newTemplateData(r)
	data.Snippets = snippets
	a.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (a *app) snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := a.snippetModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			a.serverError(w, r, err)
		}
		return
	}

	data := a.newTemplateData(r)
	data.Snippet = snippet
	a.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (a *app) snippetFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Snippet Form"))
}

func (a *app) snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 10
	snippet, err := a.snippetModel.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", snippet.ID), http.StatusSeeOther)
}
