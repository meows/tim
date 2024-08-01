package template

import (
	"github.com/a-h/templ"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/ui/template/pages"
	"github.com/timenglesf/personal-site/ui/template/partials"
)

type Pages struct {
	Base        func(title string, page templ.Component, data *shared.TemplateData) templ.Component
	Index       func(data *shared.TemplateData) templ.Component
	AdminSignup func(data *shared.TemplateData) templ.Component
	Post        func(data *shared.TemplateData) templ.Component
	CreatePost  func(data *shared.TemplateData) templ.Component
	AdminLogin  func(data *shared.TemplateData) templ.Component
}

type Partials struct {
	PageHeader  func(data *shared.TemplateData) templ.Component
	ThemeToggle func() templ.Component
	// Footer      func() templ.Component
}

func CreatePageTemplates() *Pages {
	return &Pages{
		Base:        Base,
		Index:       pages.Index,
		AdminSignup: pages.SignUpAdmin,
		AdminLogin:  pages.AdminLogin,
		Post:        pages.Post,
		CreatePost:  pages.CreatePost,
	}
}

func CreatePartialTemplates() *Partials {
	return &Partials{
		PageHeader:  partials.PageHeader,
		ThemeToggle: partials.ThemeToggle,
		//  Footer:      partials.Footer,
	}
}
