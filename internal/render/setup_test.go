package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/models"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	testApp.InProduction = false

	//log in our std lib
	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	// log.Lshortfile info about error and file
	errorLog = log.New(os.Stdout, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// we set secure to false cause here is not an https server it is testing
	session.Cookie.Secure = false
	testApp.Session = session
	// we need app not testApp in render.go
	app = &testApp
	os.Exit(m.Run())
}

//myResponseWriter is a costume response writer
type myResponseWriter struct {
}

func (w *myResponseWriter) Header() http.Header {
	var h http.Header
	return h
}

func (w *myResponseWriter) WriteHeader(i int) {

}
func (w *myResponseWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
