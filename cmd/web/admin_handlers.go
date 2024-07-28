package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/timenglesf/personal-site/internal/shared"
	"github.com/timenglesf/personal-site/internal/validator"
	"github.com/timenglesf/personal-site/ui/template"
)

func (app *application) handleDisplayAdminPage(w http.ResponseWriter, r *http.Request) {
	// If there is no admin user in the database, display the admin signup page
	adminData, err := app.user.GetAdmin()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			data := app.newAdminTemplateData(r)
			signUpPage := template.Pages.SignUpAdmin(data)
			page := template.Base("Admin Signup", false, signUpPage)
			page.Render(context.Background(), w)
			return
		}
		// TODO: Display an error message on the page using HTMX
		app.serverError(w, r, err)
	}
	app.logger.Info("Admin data", "data", adminData)
	// TODO: Display login page if not admin
	w.Write([]byte("Display the login page if not admin"))
	// TODO: Display admin dashboard if logged in as admin

	// Else display admin dashboard
}

func (app *application) handleAdminSignupPost(w http.ResponseWriter, r *http.Request) {
	// Get an admin if one exists
	adminData, err := app.user.GetAdmin()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			// TODO: Redirect to admin login page
			app.clientError(w, http.StatusBadRequest)
		}
	}
	if adminData != nil {
		// TODO: Redirect to admin login page
		w.Write([]byte("Admin already exists"))
	}

	// parse and validate form
	var form shared.AdminSignUpForm

	err = app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Email is required")
	form.CheckField(validator.ValidEmail(form.Email), "email", "Invalid email format")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "Email is too long (maximum is 100 characters)")
	form.CheckField(validator.NotBlank(form.ConfirmEmail), "confirm_email", "Confirm Email is required")
	form.CheckField(validator.EqualStrings(form.Email, form.ConfirmEmail), "confirm_email", "Emails do not match")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password is required")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "Password is too long (maximum is 100 characters)")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirm_password", "Confirm Password is required")
	form.CheckField(validator.EqualStrings(form.Password, form.ConfirmPassword), "confirm_password", "Passwords do not match")

	if !form.Valid() {
		data := app.newAdminTemplateData(r)
		data.SignUpForm = form
		page := template.Pages.SignUpAdmin(data)
		base := template.Base("Admin Signup", false, page)
		w.WriteHeader(http.StatusUnprocessableEntity)
		base.Render(context.Background(), w)
		return
	}

	// Hash password and save to database

	// Redirect to admin login page
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Admin signup successful"))
}
