package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/timenglesf/personal-site/internal/models"
)

func (app *application) handleCreateBlogPost(w http.ResponseWriter, r *http.Request) {
	title := "My very first blog post"
	content := `<h1>Welcome to My Blog</h1>
        <p>Hello, world! This is my first blog post. I'm excited to start sharing my thoughts and experiences with you. In this blog, I'll be covering a variety of topics, including technology, travel, and personal growth. Stay tuned for more updates, and thank you for joining me on this journey!</p>
        <p>Feel free to leave a comment or reach out to me through my social media channels. Let's connect and share our stories!</p>
        <footer>
            <p>Posted on July 21, 2024 by Author Name</p>
        </footer>`

	id, err := app.post.Insert(title, content, false, 0)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
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

	fmt.Fprintf(w, "%+v", post)
}
