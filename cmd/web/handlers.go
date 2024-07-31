package main

import (
	"context"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	index := app.pageTemplates.Index()
	page := app.pageTemplates.Base("Home", false, index)
	page.Render(context.Background(), w)
}
