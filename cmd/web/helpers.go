package main

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

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

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) newPostTemplateData(r *http.Request) *shared.PostTemplateData {
	return &shared.PostTemplateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) newAdminTemplateData(r *http.Request) shared.AdminTemplateData {
	return shared.AdminTemplateData{
		CurrentYear: time.Now().Year(),
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
