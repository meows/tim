package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/timenglesf/personal-site/ui/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	index := template.Pages.Index()
	page := template.Base("Home", index)
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

func (app *application) handleDisplayAdminPage(w http.ResponseWriter, r *http.Request) {
	// If not logged in display login form
	loginPage := template.Pages.Login()
	page := template.Base("Admin Login", loginPage)
	page.Render(context.Background(), w)

	// Else display admin dashboard
}
