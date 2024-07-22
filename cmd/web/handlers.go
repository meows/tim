package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/timenglesf/personal-site/ui/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	page := template.Base("Home", "Welcome to the home page.")
	page.Render(context.Background(), w)
}

func (app *application) handleGetBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	display := fmt.Sprintf("Display the blog post with ID %d", id)
	w.Write([]byte(display))
}
