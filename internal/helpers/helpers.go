package helpers

import (
	"fmt"
	"github.com/majedutd990/bookings/internal/config"
	"net/http"
	"runtime/debug"
)

//helpers we put  things in here that are useful in various parts of our application
// we call them by calling the helpers pkg

var app *config.AppConfig

// NewHelper sets up appConfig for helper
func NewHelper(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {

	//	 here we create a trace of errors using debug which is a std library
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exist := app.Session.Exists(r.Context(), "user_id")
	return exist
}
