package render

import (
	"bytes"
	"github.com/justinas/nosurf"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

// app is a reference to AppConfig that we are sending here from main
var app *config.AppConfig

//NewTemplates is the function sending the AppConfig struct here
func NewTemplates(a *config.AppConfig) {
	app = a
}

//AddDefaultData Adds some Extra data that we would love to see on every pages
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.CSRFToken = nosurf.Token(r)
	return td
}

//RenderTemplates renders html using html templates
//let us create a function that render templates
//second param is the name of the template we want to render
// instead of changing ../../template etc run the main.go using go run cmd/web/*.go
func RenderTemplates(w http.ResponseWriter, tmpl string, r *http.Request, td *models.TemplateData) {

	var tc map[string]*template.Template
	// it is a map from str to templates

	if app.UseCache {
		//this here is just for production mode where we only make the map once
		tc = app.TemplateCache
	} else {
		// when UseCache is false we create tc every time we reload a pge
		tc, _ = CreateTemplateCache()
	}

	//t is our template related to tmpl input that's been parsed
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("cannot find the corresponded template")
	}
	buf := new(bytes.Buffer)
	//adding default data to template data
	// all the tds are pointer based which means they ar addr
	td = AddDefaultData(td, r)
	//  applies parsed template to the data structure and writes result to the given writer.
	// buf makes input output better
	// here we have the parsed template we can even add more data or different data by manipulating td
	_ = t.Execute(buf, td)
	// here we write the executed template to our http.ResponseWriter
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

//CreateTemplateCache let's render all the templates and tmpl files
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	//finds all the pages and put it in a slice of strings
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		//for each page we get their names
		name := filepath.Base(page)
		// base will be for example about.page.html without any prefixes we send these
		// args from our handler
		//====================================================================================================
		//	let's create a template set
		// here we create a new template for each page in our template folder except for layouts
		// we name it by the name of the file
		//
		// we parse the file(parses given template string and return parsed template)
		// Funcs adds the elements of the argument map to the template's function map. It must be called before
		//the template is parsed. It panics if a value in the map is not a function with appropriate return type or
		//if the name cannot be used syntactically as a function in a template.
		//It is legal to overwrite elements of the map.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		//finds all the layouts
		matchesLayout, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		//if there is any puts it in ts
		if len(matchesLayout) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			//there’s also the ParseGlob function which takes glob as an argument and
			//then parses all files that matches the glob.
			// I think it merges it with our ts because we used ts.parseGlobe.
			if err != nil {
				return myCache, err
			}
		}
		//maps the current name fi (about.page.html) to its template set
		myCache[name] = ts

	}
	return myCache, nil
}
