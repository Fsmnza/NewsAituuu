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
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	s, err := app.news.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		NewsArray: s,
	})
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
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.FormValue("category")
	if title == "" || content == "" || category == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if len(title) > 20 || len(content) < 10 || len(content) > 200 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	validCategories := []string{"Students", "Staff", "Applicants", "Researches"}
	if !malika(validCategories, category) {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	id, err := app.news.Insert(title, content, category)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "News successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", id), http.StatusSeeOther)
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
		a]
