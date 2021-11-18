package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/driver"
	"github.com/majedutd990/bookings/internal/forms"
	"github.com/majedutd990/bookings/internal/helpers"
	"github.com/majedutd990/bookings/internal/models"
	"github.com/majedutd990/bookings/internal/render"
	"github.com/majedutd990/bookings/internal/repository"
	"github.com/majedutd990/bookings/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"time"
)

//Repo the repository used by handlers
var Repo *Repository

// Repository is the repo type
type Repository struct {
	App *config.AppConfig
	//	we put other things here like DB connections info
	DB repository.DataBaseRepo
}

//NewRepo creates a new repository (which basically is AppConfig in main.go)
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

//NewHandlers sets the repository for the handlers (which basically is AppConfig in main.go)
func NewHandlers(r *Repository) {
	Repo = r
}

// now we grant the access to this Repository to all these handlers

func (m *Repository) Home(writer http.ResponseWriter, r *http.Request) {

	render.Template(writer, "home.page.tmpl", r, &models.TemplateData{})
}

//let's send some data to about using a string map

func (m *Repository) About(writer http.ResponseWriter, r *http.Request) {

	render.Template(writer, "about.page.tmpl", r, &models.TemplateData{})
}

//&TemplateData{} means an empty Template Data

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation out of the session"))
		return
	}
	// we need this because we require a specific format for our dare insertion
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")
	//we put it in our string map in our template date
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)
	data := make(map[string]interface{})
	data["reservation"] = res
	var mj = &models.TemplateData{
		Form:   forms.New(nil),
		Data:   data,
		StrMap: stringMap,
	}
	render.Template(w, "make-reservation.page.tmpl", r, mj)

}

// PostReservation renders the make a reservation page and displays form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	//first thing to do when u have a form is to parse form when u have a form in it
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//we need to get our dates here than cast them in reservation object
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	rID := r.Form.Get("room_id")

	//2020-01-01 ---- m01/d02 03:04:05 pm y06-0700

	//“Mon Jan _2 15:04:05 MST 2006” ref time
	layout := "2006-01-02"
	// here is reverse we make our strings to time and date
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(rID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservations := models.Reservation{
		FirstName: r.Form.Get("firstName"),
		LastName:  r.Form.Get("lastName"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}
	//postform has all of the url values and associated data
	form := forms.New(r.PostForm)
	form.Required("firstName", "lastName", "email")
	form.MinLength("firstName", 3)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservations
		render.Template(w, "make-reservation.page.tmpl", r, &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	newReservationID, err := m.DB.InsertReservation(reservations)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	restrictions := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}
	err = m.DB.InsertRoomRestriction(restrictions)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "reservation", reservations)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the Generals page and displays form
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "generals.page.tmpl", r, &models.TemplateData{})
}

// Majors renders the Majors page and displays form
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "majors.page.tmpl", r, &models.TemplateData{})
}

// Availability renders the Availability page and displays form
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "search-availability.page.tmpl", r, &models.TemplateData{})
}

// PostAvailability renders the Availability page and displays form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	//“Mon Jan _2 15:04:05 MST 2006” ref time
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		//	no availability
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)
	render.Template(w, "choose-room.page.tmpl", r, &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomId    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJson handles requests for availability and send json responds
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)

	jResp := jsonResponse{
		Ok:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(jResp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the Contact page and displays form
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "contact.page.tmpl", r, &models.TemplateData{})
}

// ReservationSummary renders the Contact page and displays form
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	//we need to type assert it to reservation
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}
	strMap := make(map[string]string)
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	strMap["start_date"] = sd
	strMap["end_date"] = ed
	strMap["room_name"] = room.RoomName

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, "reservation-summary.page.tmpl", r, &models.TemplateData{
		Data:   data,
		StrMap: strMap,
	})
}

//ChooseRoom displays the list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	//		let's get the roomId
	// we name it id in our rout
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation out of the session"))
		return
	}
	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

//BookRoom takes urls variables make a reservation var and redirect it to make-reservation
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	var res models.Reservation
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}
