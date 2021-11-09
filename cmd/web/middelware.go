package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

////WriteToConsole next is a convention
//func WriteToConsole(next http.Handler) http.Handler {
//
//	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
//		fmt.Println("Hit The Page.")
//		//	 at the end of this we need to move on to the next
//		//	 next might be another bit of middleware
//		//   or might be a program that reads file where we actually return our mux
//		next.ServeHTTP(writer, request)
//	})
//}

// NoSurf add CSRF protection to all requests
// most middleware have a similar format to the above code
// also some may not have this anonymous function
// lets resolve CSRF Token which is going to use this package
// https://github.com/justinas/nosurf go get it
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
	//	cause it uses cookies to make sure
	//	the token it generates is available on per page basis
	//	 path "/" means entire site
	//	secure : false refers that we are not running it on a https
	//	in production we change it to true
}

//SessionLoad we should tell webserver that it has to use session using middleware
// Load and save sessions on every requests
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
