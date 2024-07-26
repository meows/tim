package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/timenglesf/personal-site/ui/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	index := template.Pages.Index()
	page := template.Base("Home", false, index)
	page.Render(context.Background(), w)
}

func (app *application) handleDisplayAdminPage(w http.ResponseWriter, r *http.Request) {
	// If there is no admin user in the database, display the admin signup page
	adminData, err := app.user.GetAdmin()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			signUpPage := template.Pages.SignUpAdmin()
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

type adminSignupForm struct {
	Email           string
	ConfirmEmail    string
	Password        string
	ConfirmPassword string
}

func (app *application) handleAdminSignupPost(w http.ResponseWriter, r *http.Request) {
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

	// Read form data
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := adminSignupForm{
		Email:           r.PostForm.Get("email"),
		ConfirmEmail:    r.PostForm.Get("confirm_email"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
	}

	fieldErrors := make(map[string]string)

	// I need to check if an email is correct format

	if strings.TrimSpace(form.Email) == "" {
		fieldErrors["email"] = "Email is required"
	} else if !isValidEmail(form.Email) {
		fieldErrors["email"] = "Invalid email format"
	} else if form.Email != form.ConfirmEmail {
		fieldErrors["confirm_email"] = "Emails do not match"
	}

	if strings.TrimSpace(form.Password) == "" {
		fieldErrors["password"] = "Password is required"
	} else if len(form.Password) < 8 {
		fieldErrors["password"] = "Password must be at least 8 characters long"
	} else if form.Password != form.ConfirmPassword {
		fieldErrors["confirm_password"] = "Passwords do not match"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprintln(w, fieldErrors)
		return
	}

	fmt.Println(form)
	w.Write([]byte("Admin signup form submitted"))
}

func (app *application) handleAdminSignupPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Admin signup page")
}

func (app *application) handleAdminLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Adming login POST")
}

func (app *application) handleAdminLoginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Admin login page")
}

func (app *application) handleAdmingLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Admin logout POST")
}
