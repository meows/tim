package main

import (
	"context"
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	index := app.pageTemplates.Index()
	page := app.pageTemplates.Base("Home", false, index)
	page.Render(context.Background(), w)
}

func (app *application) handleAdminSignupPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Admin signup page")
}

func (app *application) handleAdminLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Adming login POST")
}

func (app *application) handleAdminLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Admin logout POST")
}
