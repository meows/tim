package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/internal/validator"
	"github.com/timenglesf/personal-site/ui/template"
	"github.com/yuin/goldmark"
)

// Create Post Handlers

// Displays form to create a blog post
func (app *application) handleDisplayCreatePostForm(w http.ResponseWriter, r *http.Request) {
	// TODO: Make sure user is logged in and is an admin

	page := template.Pages.CreatePost(shared.AdminTemplateData{})
	w.Header().Set("Content-Type", "text/html")
	base := template.Base("Create Post", true, page)
	if err := base.Render(context.Background(), w); err != nil {
		app.serverError(w, r, err)
		return
	}
}

// Saves a newly created blog to db and redirects to view the post
func (app *application) handleCreateBlogPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := shared.BlogPostFormData{
		Title:   r.PostFormValue("title"),
		Content: r.PostFormValue("content"),
	}

	// Validate form data
	form.CheckField(validator.NotBlank(form.Title), "title", "Title is required")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field is too long (maximum is 100 characters)")
	form.CheckField(validator.NotBlank(form.Content), "content", "Content is required")

	if !form.Valid() {
		fmt.Println("Form has errors")
		fmt.Println(form.Content)
		data := app.newAdminTemplateData(r)
		data.BlogForm = form
		page := template.Pages.CreatePost(data)
		base := template.Base("Create Post", true, page)
		w.WriteHeader(http.StatusUnprocessableEntity)
		base.Render(context.Background(), w)
		return
	}

	id, err := app.post.Insert(form.Title, form.Content, false, 0)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

// View blog Post by ID
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
