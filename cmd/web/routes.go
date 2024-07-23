package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /blog/{id}", dynamic.ThenFunc(app.handleGetBlogPost))

	// Admin routes
	mux.Handle("GET /_", dynamic.ThenFunc(app.handleDisplayAdminPage))
	mux.Handle("GET /_/signup", dynamic.ThenFunc(app.handleAdminSignupPage))
	mux.Handle("POST /_/signup", dynamic.ThenFunc(app.handleAdminSignupPost))
	mux.Handle("GET /_/login", dynamic.ThenFunc(app.handleAdminLoginPage))
	mux.Handle("POST /_/login", dynamic.ThenFunc(app.handleAdminLoginPost))
	mux.Handle("POST /_/logout", dynamic.ThenFunc(app.handleAdmingLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
