package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/majedutd990/bookings/pkg/config"
	"github.com/majedutd990/bookings/pkg/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	//// ==========================we used pat package
	////	 in here we create a new http handler which we often call Mux or a multiplexer
	////	we use an external package for routing
	////	we use bmizerany rout package
	////	remember mux is a htt handler so
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
	return mux
}
