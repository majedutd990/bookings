package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	//// ==========================we used pat package
	////	 in here we create a new http handler which we often call Mux or a multiplexer
	////	we use an external package for routing
	////	we use bmizerany rout package
	////	remember mux is a http handler so
	//mux := pat.New()
	//
	//// here we set up our routs
	//// we use its get function to set up our routes
	//mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	//mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	//========================= here we use chi =========
	mux := chi.NewRouter()
	// here we can also use the advantage of middlewares
	// they apparently must come before our routes
	// for example here we use Recoverer:Gracefully absorb panics and prints the stack trace
	// should use it without parenthesis: maybe we send a function ref here
	mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	// my middleware example

	//CSRF attack middleWare
	mux.Use(NoSurf)
	// use session all the time
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suites", handlers.Repo.Majors)
	//searching for availability we need both post and get here the only difference will be the handler
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJson)
	//chose the room
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	// book-room
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.LogOut)
	//============= static files=============
	//we have to tell this router how to return our static files
	// we have to create a file server a place that go gets these file from
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	// here we define routes that need to be protected or only be shown to registered users
	// all the routes will be /admin/dashboard for example
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
	})
	return mux
}
