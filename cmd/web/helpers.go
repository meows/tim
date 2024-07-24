package main

import (
	"net/http"
	"regexp"
	"runtime/debug"
	"time"

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

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func isValidEmail(email string) bool {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegexPattern)

	return re.MatchString(email)
}

func (app *application) newPostTemplateData(r *http.Request) *shared.PostTemplateData {
	return &shared.PostTemplateData{
		CurrentYear: time.Now().Year(),
	}
}
