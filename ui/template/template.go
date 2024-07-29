package template

import (
	"github.com/a-h/templ"
	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/ui/template/pages"
	"github.com/timenglesf/personal-site/ui/template/partials"
)

type PagesStruct struct {
	Index       func() templ.Component
	SignUpAdmin func(data shared.AdminTemplateData) templ.Component
	Post        func(data shared.PostTemplateData) templ.Component
	CreatePost  func(data shared.AdminTemplateData) templ.Component
}

type PartialStruct struct {
	PageHeader  func(isAdmin bool) templ.Component
	ThemeToggle func() templ.Component
}

// TODO: Put these on the app struct
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
