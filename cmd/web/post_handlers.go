package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/ui/template"
	"github.com/yuin/goldmark"
)

// Create Post Handlers

func (app *application) handleDisplayCreatePostForm(w http.ResponseWriter, r *http.Request) {
	// TODO: Make sure user is logged in and is an admin

	page := template.Pages.CreatePost()
	w.Header().Set("Content-Type", "text/html")
	base := template.Base("Create Post", true, page)
	if err := base.Render(context.Background(), w); err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) handleCreateBlogPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostFormValue("title")
	content := r.PostFormValue("content")

	id, err := app.post.Insert(title, content, false, 0)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func (app *application) handleGetBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	post, err := app.post.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newPostTemplateData(r)
	data.BlogPost = post

	// convert markdown to html
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(data.BlogPost.Content), &buf); err != nil {
		app.serverError(w, r, err)
		return
	}

	data.BlogPost.ContentHTML = buf.String()
	page := template.Pages.Post(*data)

	w.Header().Set("Content-Type", "text/html")

	base := template.Base(data.BlogPost.Title, false, page)
	base.Render(context.Background(), w)
}

func (app *application) handleGetLatestBlogPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := app.post.Latest(false)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	for _, post := range posts {
		fmt.Fprintf(w, "%+v\n", post)
	}
}
