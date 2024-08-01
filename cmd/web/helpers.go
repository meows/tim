package main

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/form/v4"
	"github.com/timenglesf/personal-site/internal/shared"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) logServerError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
}

func (app *application) logServerWarning(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Warn(err.Error(), "method", method, "uri", uri, "trace", trace)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// func (app *application) newPostTemplateData(r *http.Request) *shared.PostTemplateData {
// 	return &shared.PostTemplateData{
// 		CurrentYear: time.Now().Year(),
// 	}
// }

func (app *application) newTemplateData(r *http.Request) shared.TemplateData {
	return shared.TemplateData{
		IsAuthenticated: app.isAuthenticated(r),
		IsAdmin:         app.isAdmin(r),
		CurrentYear:     time.Now().Year(),
	}
}

// a helper function to decode form data into a struct using the automatic form decoder
func (app *application) decodeForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	err := app.formDecoder.Decode(&dst, r.PostForm)

	var inValidDecoderError *form.InvalidDecoderError

	if err != nil {
		if errors.As(err, &inValidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}

// returns true if the session data contains the key "authenticatedUserID"
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

// returns true if the session data contains the key "isAdminRole"
func (app *application) isAdmin(r *http.Request) bool {
	return app.sessionManager.GetBool(r.Context(), "isAdminRole")
}

func (app *application) renderPage(w http.ResponseWriter, r *http.Request, templateFunc func(data *shared.TemplateData) templ.Component, title string, data *shared.TemplateData) {
	page := templateFunc(data)
	base := app.pageTemplates.Base(title, page, data)
	err := base.Render(r.Context(), w)
	if err != nil {
		app.serverError(w, r, err)
	}
}
