package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/justinas/nosurf"
	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/internal/validator"
	"github.com/yuin/goldmark"
)

// Create Post Handlers

// Displays form to create a blog post
func (app *application) handleDisplayCreatePostForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	if !data.IsAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.renderPage(w, r, app.pageTemplates.CreatePost, "Create Post", &data)
}

// Saves a newly created blog to db and redirects to view the post
func (app *application) handleCreateBlogPost(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Creating new blog post")
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
		app.renderPage(w, r, app.pageTemplates.CreatePost, "Create Post", &data)
		return
	}

	// Get author id
	userId := app.sessionManager.GetString(r.Context(), sessionUserId)

	id, err := app.post.Insert(form.Title, form.Content, false, userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flashSuccess", "Post succesfully created!")

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

// Render blog Post by Title
func (app *application) handleGetBlogPost(w http.ResponseWriter, r *http.Request) {
	titleId := r.PathValue("slug")
	targetPostTitle, err := url.QueryUnescape(titleId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	post, err := app.post.GetPostByTitle(targetPostTitle)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Reject unauthorized access to private posts
	if post.Private {
		if !app.isAdmin(r) {
			app.logger.Warn("unauthorized access to url", "url", r.URL.Path, "ip", r.RemoteAddr)
			referer := r.Referer()
			if referer != "" {
				http.Redirect(w, r, referer, http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}
	}

	data := app.newTemplateData(r)
	data.BlogPost = post

	// convert markdown to html
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(data.BlogPost.Content), &buf); err != nil {
		app.logger.Error("Error converting markdown to html", "error", err)
		app.serverError(w, r, err)
		return
	}

	data.BlogPost.Content = buf.String()

	// Get flash message from session
	flashSuccess := app.sessionManager.PopString(r.Context(), "flashSuccess")
	if flashSuccess != "" {
		data.Flash = &shared.FlashMessage{Message: flashSuccess, Type: "Post created successfully"}
	}

	app.renderBlogPostPage(w, r, post.Title, &data)
	// app.renderPage(w, r, app.pageTemplates.Post, "Post", &data)
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

func (app *application) handleBlogPostUpdate(w http.ResponseWriter, r *http.Request) {
	sentCSRFTOKEN, err := r.Cookie("csrf_token")
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if !nosurf.VerifyToken(nosurf.Token(r), sentCSRFTOKEN.Value) {
		return
	}

	slug := r.PathValue("slug")
	query := r.URL.Query()

	title, _ := url.QueryUnescape(slug)

	post, err := app.post.GetPostByTitle(title)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	for key, value := range query {
		switch key {
		case "title":
			post.Title = value[0]
		case "content":
			post.Content = value[0]
		case "private":
			post.Private = !post.Private
		}
	}

	if err := app.post.Update(post); err != nil {
		app.serverError(w, r, err)
		return
	}

	// Send updated blog post row if this is an updated to the private column
	if query.Get("private") != "" {
		updatedRowComponenet := app.partialTemplates.DashboardBlogPostRow(post)
		if err = updatedRowComponenet.Render(r.Context(), w); err != nil {
			app.serverError(w, r, err)
		}
	}
}

func (app *application) handleDisplayEditPostForm(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	id, err := strconv.Atoi(slug)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	post, err := app.post.GetPostByID(uint(id))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var form shared.BlogPostFormData
	form.Title = post.Title
	form.Content = post.Content

	templateData := app.newTemplateData(r)
	templateData.BlogForm = form
	templateData.BlogPost = post

	page := app.pageTemplates.EditPost(&templateData)
	if err = page.Render(r.Context(), w); err != nil {
		fmt.Println(err)
		app.serverError(w, r, err)
	}
}

func (app *application) handleBlogPostEdit(w http.ResponseWriter, r *http.Request) {
	var form struct {
		shared.BlogPostFormData
		ID uint `form:"id"`
	}
	if err := app.decodeForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var blogFormData shared.BlogPostFormData
	blogFormData.Title = form.Title
	blogFormData.Content = form.Content

	// Validate form data
	blogFormData.CheckField(validator.NotBlank(form.Title), "title", "Title is required")
	blogFormData.CheckField(validator.MaxChars(form.Title, 100), "title", "This field is too long (maximum is 100 characters)")
	blogFormData.CheckField(validator.NotBlank(form.Content), "content", "Content is required")

	data := app.newTemplateData(r)
	data.BlogForm = blogFormData

	if !blogFormData.Valid() {
		data.BlogPost = &models.Post{}
		data.BlogPost.ID = form.ID
		page := app.pageTemplates.EditPost(&data)
		if err := page.Render(r.Context(), w); err != nil {
			app.serverError(w, r, err)
		}
		return
	}

	// TODO: Update post in db and redirect to view post
}

func printTemplateData(templateData interface{}) {
	v := reflect.ValueOf(templateData)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		fmt.Printf("Key: %s, Value: %v\n", field.Name, value)
	}
}
