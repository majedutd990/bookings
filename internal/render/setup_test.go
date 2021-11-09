package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	testApp.InProduction = false

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
