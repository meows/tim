package main

import "github.com/timenglesf/personal-site/internal/models"

type adminTemplateData struct {
	CurrentYear int
	Blogs       any
	Forms       any
	Admin       models.User
}

type templateData struct {
	CurrentYear int
	Blogs       any
	Blog        any
	Form        any
	Flash       string
	CSRFtoken   string
	User        models.User
}
