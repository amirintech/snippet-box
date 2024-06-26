package main

import (
	"errors"
	"fmt"
	"github.com/amirintech/snippet-box/internal/models"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
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
	data := a.newTemplateData(r)
	a.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (a *app) snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	// parse & extract form data
	if err := r.ParseForm(); err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	expires, err := strconv.Atoi(r.FormValue("expires"))
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	// validate form fields
	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "Title cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "Title cannot be more than 100 characters"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "Content cannot be blank"
	} else if utf8.RuneCountInString(content) > 50_000 {
		fieldErrors["content"] = "Content cannot be more than 50,000 characters"
	}

	const (
		day  = 1
		week = 7
		year = 365
	)
	if expires != day && expires != week && expires != year {
		fieldErrors["expires"] = "Expires value must be 1, 7, or 365"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	// create snippet
	snippet, err := a.snippetModel.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", snippet.ID), http.StatusSeeOther)
}
