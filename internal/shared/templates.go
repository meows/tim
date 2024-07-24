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
	Blogs       any
	Forms       any
	Admin       models.User
}

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
