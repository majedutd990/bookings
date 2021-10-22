package handlers

import (
	"github.com/majedutd990/bookings/pkg/config"
	"github.com/majedutd990/bookings/pkg/models"
	"github.com/majedutd990/bookings/pkg/render"
	"net/http"
)

//Repo the repository used by handlers
var Repo *Repository

// Repository is the repo type
type Repository struct {
	App *config.AppConfig
	//	we put other things here like DB connections info
}

//NewRepo creates a new repository (which basically is AppConfig in main.go)
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the repository for the handlers (which basically is AppConfig in main.go)
func NewHandlers(r *Repository) {
	Repo = r
}

// now we grant the access to this Repository to all these handlers

func (m *Repository) Home(writer http.ResponseWriter, r *http.Request) {
	//let's get the remote addr
	remoteIP := r.RemoteAddr
	//put it in our session
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplates(writer, "home.page.tmpl", &models.TemplateData{})
}

//let's send some data to about using a string map

func (m *Repository) About(writer http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, Again!"
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	render.RenderTemplates(writer, "about.page.tmpl", &models.TemplateData{
		StrMap: stringMap,
	})
}

//&TemplateData{} means an empty Template Data

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w,"make-reservation.page.tmpl",&models.TemplateData{})
}

// Generals renders the Generals page and displays form
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w,"generals.page.tmpl",&models.TemplateData{})
}

// Majors renders the Majors page and displays form
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w,"majors.page.tmpl",&models.TemplateData{})
}

// Availability renders the Availability page and displays form
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w,"search-availability.page.tmpl",&models.TemplateData{})
}
// Contact renders the Contact page and displays form
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w,"contact.page.tmpl",&models.TemplateData{})
}