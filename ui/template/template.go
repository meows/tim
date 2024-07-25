package template

import (
	"github.com/a-h/templ"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/ui/template/pages"
	"github.com/timenglesf/personal-site/ui/template/partials"
)

type PagesStruct struct {
	Index       func() templ.Component
	SignUpAdmin func() templ.Component
	Post        func(shared.PostTemplateData) templ.Component
	CreatePost  func(shared.AdminTemplateData) templ.Component
}

type PartialStruct struct {
	PageHeader  func(isAdmin bool) templ.Component
	ThemeToggle func() templ.Component
}

var Pages = PagesStruct{
	Index:       pages.Index,
	SignUpAdmin: pages.SignUpAdmin,
	Post:        pages.Post,
	CreatePost:  pages.CreatePost,
}

var Partials = PartialStruct{
	PageHeader:  partials.PageHeader,
	ThemeToggle: partials.ThemeToggle,
}
