package template

import (
	"github.com/a-h/templ"
	"github.com/timenglesf/personal-site/ui/template/pages"
	"github.com/timenglesf/personal-site/ui/template/partials"
)

type PagesStruct struct {
	Index func() templ.Component
	Login func() templ.Component
}

type PartialStruct struct {
	PageHeader  func(themeToggle templ.Component) templ.Component
	ThemeToggle func() templ.Component
}

var Pages = PagesStruct{
	Index: pages.Index,
	Login: pages.Login,
}

var Partials = PartialStruct{
	PageHeader:  partials.PageHeader,
	ThemeToggle: partials.ThemeToggle,
}
