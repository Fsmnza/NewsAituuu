package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/news", app.showNews)
	mux.HandleFunc("/ne", app.showDepartments)
	mux.HandleFunc("/news/create", app.createNews)
	mux.HandleFunc("/news/departments", app.createDepartments)
	mux.HandleFunc("/news/creationPage", app.creationPage)
	mux.HandleFunc("/departments/creatDep", app.creationPg)
	mux.HandleFunc("/news/contacts", app.contacts)
	mux.HandleFunc("/news/students", app.filterCategory("Students"))
	mux.HandleFunc("/news/applicants", app.filterCategory("Applicants"))
	mux.HandleFunc("/news/researches", app.filterCategory("Researches"))
	mux.HandleFunc("/news/staff", app.filterCategory("Staff"))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
