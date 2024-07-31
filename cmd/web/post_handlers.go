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
	"github.com/yuin/goldmark"
)

// Create Post Handlers

// Displays form to create a blog post
func (app *application) handleDisplayCreatePostForm(w http.ResponseWriter, r *http.Request) {
	// TODO: Make sure user is logged in and is an admin

	data := app.newTemplateData(r)
	app.renderAdminPage(w, r, app.pageTemplates.CreatePost, "Create Post", &data)
	// page := app.pageTemplates.CreatePost(&data)
	//	w.Header().Set("Content-Type", "text/html")
	// base := template.Base("Create Post", true, page)
	// if err := base.Render(context.Background(), w); err != nil {
	// app.serverError(w, r, err)
	// return
	// }
}

// Saves a newly created blog to db and redirects to view the post
func (app *application) handleCreateBlogPost(w http.ResponseWriter, r *http.Request) {
	var form shared.BlogPostFormData

	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate form data
	form.CheckField(validator.NotBlank(form.Title), "title", "Title is required")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field is too long (maximum is 100 characters)")
	form.CheckField(validator.NotBlank(form.Content), "content", "Content is required")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.BlogForm = form
		page := app.pageTemplates.CreatePost(&data)
		base := app.pageTemplates.Base("Create Post", true, page)
		w.WriteHeader(http.StatusUnprocessableEntity)
		base.Render(context.Background(), w)
		app.renderAdminPage(w, r, app.pageTemplates.CreatePost, "Create Post", &data)
		return
	}

	id, err := app.post.Insert(form.Title, form.Content, false, 0)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flashSuccess", "Post succesfully created!")

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

// View blog Post by ID
func (app *application) handleGetBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("slug"))
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

	data := app.newTemplateData(r)
	data.BlogPost = post

	// convert markdown to html
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(data.BlogPost.Content), &buf); err != nil {
		app.serverError(w, r, err)
		return
	}

	data.BlogPost.ContentHTML = buf.String()

	// Get flash message from session
	flashSuccess := app.sessionManager.PopString(r.Context(), "flashSuccess")
	if flashSuccess != "" {
		data.Flash = shared.FlashMessage{Message: flashSuccess, Type: "Post created successfully"}
	}

	//	page := app.pageTemplates.Post(*data)
	//
	//	w.WriteHeader(http.StatusCreated)
	//
	//	base := template.Base(data.BlogPost.Title, false, page)
	//
	//	base.Render(r.Context(), w)
	app.renderAdminPage(w, r, app.pageTemplates.Post, "Post", &data)
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
