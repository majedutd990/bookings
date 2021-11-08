package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/handlers"
	"github.com/majedutd990/bookings/internal/models"
	"github.com/majedutd990/bookings/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

// we made this app var package level, so we have access to it in middleware file
var app config.AppConfig

// session as above
var session *scs.SessionManager

func main() {

	//what we put in the session
	gob.Register(models.Reservation{})
	// change this to true when in production
	app.InProduction = false

	// let declare our sessions
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	//set the life cycle of the session
	// by default it uses cookies to store them however we can use some dbs like badger store , sqlite , ...
	session.Cookie.Persist = true
	// session persists after closing the browser
	session.Cookie.SameSite = http.SameSiteLaxMode
	// how strict do u wan to be about what site this cookie applies to
	session.Cookie.Secure = app.InProduction
	// insists that cookie be encrypted related to https in here we or development mode
	// we use false in production we have to do the

	//set our config session, so it will be visible to other pkgs like handlers
	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	// the below code means we are in dev mode
	app.UseCache = app.InProduction
	//because it's type is pointer to AppConfig, here we send a reference
	render.NewTemplates(&app)

	// now let's create repo for our handlers and pass our app config here
	repo := handlers.NewRepo(&app)
	// now we pass it to handlers which uses it to make Repo Var
	handlers.NewHandlers(repo)

	// Repo is the repository used by handlers we make it above
	// use it below cause home and about are of type repository function

	// - obsolete
	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	//fmt.Println(fmt.Sprintf("Starting Application On Port %s.", portNumber))
	//_ = http.ListenAndServe(portNumber, nil)

	//	 new version using pat package
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	fmt.Println(fmt.Sprintf("Starting Application On Port %s.", portNumber))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
