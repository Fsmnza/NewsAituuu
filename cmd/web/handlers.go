package main

import (
	"alexedwards.net/snippetbox/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	_ "strings"
	_ "unicode/utf8"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//userID, ok := app.session.Get(r, "authenticatedUserID").(int)
	//if !ok {
	//	app.notFound(w)
	//	return
	//}
	//isAdmin, err := app.users.IsAdmin(userID)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}

	s, err := app.news.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	td := &templateData{
		NewsArray: s,
		UserRole:  "user",
	}

	//if isAdmin {
	//	td.UserRole = "admin"
	//}

	app.render(w, r, "home.page.tmpl", td)
}

func (app *application) showNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	n, err := app.news.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	if n == nil {
		app.notFound(w)
		return
	}
	app.render(w, r, "show.page.tmpl", &templateData{
		News: n,
	})
}
func (app *application) creationPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{})
}

func (app *application) createNews(w http.ResponseWriter, r *http.Request) {
	userID, ok := app.session.Get(r, "authenticatedUserID").(int)
	if !ok {
		app.notFound(w)
		return
	}
	isTeacher, err := app.users.IsTeacher(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if !isTeacher {
		app.clientError(w, http.StatusForbidden)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	app.session.Put(r, "flash", "News successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d"), http.StatusSeeOther)
}
func malika(slice []string, s string) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}
func (app *application) filterCategory(category string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.renderCategoryPage(w, r, category)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
}
func (app *application) renderCategoryPage(w http.ResponseWriter, r *http.Request, category string) error {
	newsArray, err := app.news.GetByCategory(category)
	if err != nil {
		return err
	}
	app.render(w, r, "category.page.tmpl", &templateData{
		Category:  category,
		NewsArray: newsArray,
	})
	return nil
}

func (app *application) contacts(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "contacts.page.tmpl", &templateData{})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Println(name + email + password)
	err = app.users.Insert(name, email, password)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Sign up is successful")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	ID, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.session.Put(r, "flash", "Invalid email or password.")
			app.render(w, r, "login.page.tmpl", nil)
			return
		}
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "authenticatedUserID", ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//func (app *application) admin(w http.ResponseWriter, r *http.Request) {
//	userID, ok := app.session.Get(r, "authenticatedUserID").(int)
//	if !ok {
//		app.notFound(w)
//		return
//	}
//	role, err := app.users.GetRole(userID)
//	if err != nil {
//		app.serverError(w, err)
//		return
//	}
//	if role != "admin" {
//		app.clientError(w, http.StatusForbidden)
//		return
//	}
//	app.render(w, r, "admin.page.tmpl", &templateData{})
//}
