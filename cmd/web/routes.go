package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))

	// Posts
	mux.Handle("GET /posts/view/{slug}", dynamic.ThenFunc(app.handleGetBlogPost))
	mux.Handle("GET /posts/latest", dynamic.ThenFunc(app.handleGetLatestBlogPosts))

	adminProtected := dynamic.Append(app.requireAdmin)
	// Admin routes
	mux.Handle("POST /posts/create", dynamic.ThenFunc(app.handleCreateBlogPost))
	mux.Handle("GET /posts/create", adminProtected.ThenFunc(app.handleDisplayCreatePostForm))

	mux.Handle("GET /admin/{$}", dynamic.ThenFunc(app.handleDisplayAdminPage))
	mux.Handle("GET /admin/signup", dynamic.ThenFunc(app.handleAdminSignupPage))
	mux.Handle("POST /admin/signup", dynamic.ThenFunc(app.handleAdminSignupPost))
	mux.Handle("GET /admin/login", dynamic.ThenFunc(app.handleAdminLoginPage))
	mux.Handle("POST /admin/login", dynamic.ThenFunc(app.handleAdminLoginPost))
	mux.Handle("POST /admin/logout", dynamic.ThenFunc(app.handleAdminLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
