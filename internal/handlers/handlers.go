package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/models"
	"github.com/majedutd990/bookings/internal/render"
	"log"
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
	render.RenderTemplates(writer, "home.page.tmpl", r, &models.TemplateData{})
}

//let's send some data to about using a string map

func (m *Repository) About(writer http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, Again!"
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	render.RenderTemplates(writer, "about.page.tmpl", r, &models.TemplateData{
		StrMap: stringMap,
	})
}

//&TemplateData{} means an empty Template Data

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "make-reservation.page.tmpl", r, &models.TemplateData{})
}

// Generals renders the Generals page and displays form
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "generals.page.tmpl", r, &models.TemplateData{})
}

// Majors renders the Majors page and displays form
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "majors.page.tmpl", r, &models.TemplateData{})
}

// Availability renders the Availability page and displays form
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "search-availability.page.tmpl", r, &models.TemplateData{})
}

// PostAvailability renders the Availability page and displays form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("The start date is %s and the end date is %s.\n", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJson handles requests for availability and send json responds
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {

	jresp := jsonResponse{
		Ok:      true,
		Message: "Available!",
	}
	out, err := json.MarshalIndent(jresp, "", "     ")
	if err != nil {
		log.Fatal(err)

	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the Contact page and displays form
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "contact.page.tmpl", r, &models.TemplateData{})
}
