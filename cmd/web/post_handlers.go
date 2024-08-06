package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"

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
		fmt.Println(err)
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

	id, err := app.post.Insert(form.Title, form.Content, false, 0)
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
			fmt.Println(referer)
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
	fmt.Println(data.BlogPost)

	// convert markdown to html
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(data.BlogPost.Content), &buf); err != nil {
		app.logger.Error("Error converting markdown to html", "error", err)
		app.serverError(w, r, err)
		return
	}

	data.BlogPost.ContentHTML = buf.String()

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
		updatedRowComponenet := app.partialTemplates.DashboardBlogPostRow(&post)
		updatedRowComponenet.Render(r.Context(), w)
	}
}
