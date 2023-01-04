package main

import (
	"fmt"
	"github.com/calvarado2004/go-concurrency/data"
	"html/template"
	"net/http"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float32
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	User          *data.User
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, templates string, templateData *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}

	var templateSlice []string

	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", pathToTemplates, templates))

	for _, currentIteration := range partials {
		templateSlice = append(templateSlice, currentIteration)
	}

	if templateData == nil {
		templateData = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, app.AddDefaultData(templateData, r)); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {

	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Authenticated = app.IsAuthenticated(r)
	if td.Authenticated {

		// get more user information
		user, ok := app.Session.Get(r.Context(), "user").(data.User)
		if !ok {
			app.ErrorLog.Println("Cannot get user from session")
			return td
		} else {
			td.User = &user
		}

	}

	td.Now = time.Now()

	return td
}

func (app *Config) IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "userID")

	return exists
}
