package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	posts, err := app.post.Latest(false)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.BlogPosts = posts
	app.renderPage(w, r, app.pageTemplates.Index, "Home", &data)
}
