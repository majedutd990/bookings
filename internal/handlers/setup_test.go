package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/models"
	"github.com/majedutd990/bookings/internal/render"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var functions = template.FuncMap{}
var pathToTemplate = "./../../templates"

func getRoutes() http.Handler {

	gob.Register(models.Reservation{})
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	//we set UseCache to false otherwise it will use the original create template cache in render package
	app.UseCache = true
	render.NewTemplates(&app)

	repo := NewRepo(&app)
	NewHandlers(repo)
	//route file
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//we set csrf token off for the test
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suites", Repo.Majors)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJson)

	mux.Get("/contact", Repo.Contact)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler

}
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matchesLayout, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
		if err != nil {
			return myCache, err
		}
		if len(matchesLayout) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts

	}
	return myCache, nil
}
