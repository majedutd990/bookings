package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/driver"
	"github.com/majedutd990/bookings/internal/handlers"
	"github.com/majedutd990/bookings/internal/helpers"
	"github.com/majedutd990/bookings/internal/models"
	"github.com/majedutd990/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

// we made this app var package level, so we have access to it in middleware file
var app config.AppConfig

// session as above
var session *scs.SessionManager

// setting up log vars
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)
	log.Println("Starting mail listener....")
	listenForMail()
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
func run() (*driver.DB, error) {

	//what we put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	//let's read some of our configuration atr values from command line as a flag with -
	// change this to true when in production
	inProduction := flag.Bool("production", true, "Application is in production!")
	useCache := flag.Bool("cache", true, "Use template cache (in production true)!")
	// let's specify database info because we don't want to put them in our code base we get it from command line
	dbName := flag.String("dbname", "", "DataBase Name")
	dbHost := flag.String("dbhost", "localhost", "DataBase Host")
	dbUser := flag.String("dbuser", "", "DataBase username")
	dbPass := flag.String("dbpass", "", "DataBase password")
	dbPort := flag.String("dbport", "5432", "DataBase port")
	dbSSL := flag.String("dbssl", "disable", "DataBase SSL settings (disable, prefer,require)")
	// let's parse them
	flag.Parse()
	// prefer use ssl if it exists, require u must have it
	// remember all of them are pointers so should use * indirect access operator
	// we set all of them in run.sh
	// also we need to check for the necessary flags
	// we can also use .env file and also a package for it go check it out

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags!")
		os.Exit(1)
	}

	app.InProduction = *inProduction
	//log in our std lib
	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	// log.Lshortfile info about error and file
	errorLog = log.New(os.Stdout, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
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
	// insists that cookie be encrypted related to https in here we are in development mode
	// we use false in production.
	//set our config session, so it will be visible to other pkgs like handlers
	app.Session = session

	//connect to database
	log.Println("Connecting to database")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)

	db, err := driver.ConnectSql(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database. Dying!")
	}
	log.Println("Connected to database")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("cannot create template cache")
		return db, err
	}
	app.TemplateCache = tc
	// the below code means we are in dev mode
	app.UseCache = *useCache
	//because it's type is pointer to AppConfig, here we send a reference
	render.NewRenderer(&app)

	// now let's create repo for our handlers and pass our app config here
	repo := handlers.NewRepo(&app, db)
	// now we pass it to handlers which uses it to make Repo Var
	handlers.NewHandlers(repo)

	// Repo is the repository used by handlers we make it above
	// use it below cause home and about are of type repository function

	//we can pass app to helpers here
	helpers.NewHelper(&app)
	// - obsolete
	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	//fmt.Println(fmt.Sprintf("Starting Application On Port %s.", portNumber))
	//_ = http.ListenAndServe(portNumber, nil)

	//	 new version using pat package
	return db, nil
}
