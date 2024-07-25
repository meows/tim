package shared

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/timenglesf/personal-site/internal/models"
)

const DateLayout = "January 2, 2006"

type PostTemplateData struct {
	CurrentYear int
	// Posts       []models.Post
	BlogPost models.Post
	Form     any
	// Flash       string
	// CSRFtoken   string
	//  User        models.User
}

type AdminTemplateData struct {
	CurrentYear int
	BlogPost    models.Post
	// Posts       []models.Post
	Forms    any
	Admin    models.User
	Form     any
	BlogForm BlogPostFormData
}

type BlogPostFormData struct {
	Title       string
	Content     string
	FieldErrors map[string]string
	// Tags
}

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
