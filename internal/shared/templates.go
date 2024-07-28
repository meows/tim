package shared

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/timenglesf/personal-site/internal/models"
	"github.com/timenglesf/personal-site/internal/validator"
)

const DateLayout = "January 2, 2006"

type PostTemplateData struct {
	// Posts       []models.Post
	BlogPost    models.Post
	Flash       string
	CurrentYear int
	Form        any
	// CSRFtoken   string
	//  User        models.User
}

type AdminTemplateData struct {
	BlogPost    models.Post
	Admin       models.User
	BlogForm    BlogPostFormData
	SignUpForm  AdminSignUpForm
	CurrentYear int
	// Posts       []models.Post
}

type AdminSignUpForm struct {
	Email               string `form:"email"`
	ConfirmEmail        string `form:"confirm_email"`
	Password            string `form:"password"`
	ConfirmPassword     string `form:"confirm_password"`
	validator.Validator `form:"-"`
}

type BlogPostFormData struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	validator.Validator `form:"-"`
	// Tags
}

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
