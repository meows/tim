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

func (app *application) handleCreateBlogPost(w http.ResponseWriter, r *http.Request) {
	title := "My second blog post"
	content := `<h1>The Second Blog Post</h1>
        <p>Hello, world! This is my second blog post. I'm excited to start sharing my thoughts and experiences with you. In this blog, I'll be covering a variety of topics, including technology, travel, and personal growth. Stay tuned for more updates, and thank you for joining me on this journey!</p>
        <p>Feel free to leave a comment or reach out to me through my social media channels. Let's connect and share our stories!</p>
        <footer>
            <p>Posted on July 21, 2024 by Author Name</p>
        </footer>`

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
