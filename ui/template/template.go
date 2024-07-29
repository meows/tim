package template

import (
	"github.com/a-h/templ"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/ui/template/pages"
	"github.com/timenglesf/personal-site/ui/template/partials"
)

type Pages struct {
	Base        func(title string, isAdmin bool, page templ.Component) templ.Component
	Index       func() templ.Component
	AdminSignup func(data shared.AdminTemplateData) templ.Component
	Post        func(data shared.PostTemplateData) templ.Component
	CreatePost  func(data shared.AdminTemplateData) templ.Component
	AdminLogin  func(data shared.AdminTemplateData) templ.Component
}

type Partials struct {
	PageHeader  func(isAdmin bool) templ.Component
	ThemeToggle func() templ.Component
	// Footer      func() templ.Component
}

func CreatePageTemplates() *Pages {
	return &Pages{
		Base:        Base,
		Index:       pages.Index,
		AdminSignup: pages.SignUpAdmin,
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
