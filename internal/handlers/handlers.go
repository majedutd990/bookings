package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/driver"
	"github.com/majedutd990/bookings/internal/forms"
	"github.com/majedutd990/bookings/internal/models"
	"github.com/majedutd990/bookings/internal/render"
	"github.com/majedutd990/bookings/internal/repository"
	"github.com/majedutd990/bookings/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"strings"
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

//NewTestRepo creates a new repository (which basically is AppConfig in main.go)
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
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
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// we need this because we require a specific format for our dare insertion
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(rID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	room, err := m.DB.GetRoomByID(roomID)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "no such room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservations.Room.RoomName = room.RoomName

	//postform has all of the url values and associated data
	form := forms.New(r.PostForm)
	form.Required("firstName", "lastName", "email")
	form.MinLength("firstName", 3)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservations
		http.Error(w, "invalid form", http.StatusSeeOther)
		render.Template(w, "make-reservation.page.tmpl", r, &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	newReservationID, err := m.DB.InsertReservation(reservations)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation to data base")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// send notification mails to user
	htmlMsg := fmt.Sprintf(`
	<strong>Reservation Confirmation</strong>
	Dear %s:,</br>
	This is to confirm your reservation from %s to %s`,
		reservations.FirstName,
		reservations.StartDate.Format("2006-01-02"),
		reservations.StartDate.Format("2006-01-02"))
	mail := models.MailData{
		To:       reservations.Email,
		From:     "majedutd@gmail.com",
		Subject:  "Reservation Confirmation!",
		Content:  htmlMsg,
		Template: "basic.html",
	}
	m.App.MailChan <- mail
	// send notification mails to owner
	htmlMsg = fmt.Sprintf(`
	<strong>Reservation Confirmation</strong> </br>
	A Reservation has been made for %s from %s to %s`,
		reservations.FirstName,
		reservations.StartDate.Format("2006-01-02"),
		reservations.StartDate.Format("2006-01-02"))
	mail = models.MailData{
		To:      "majedutd@gmail.com",
		From:    "majedutd@gmail.com",
		Subject: "Reservation Notification!",
		Content: htmlMsg,
	}
	m.App.MailChan <- mail
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
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	//“Mon Jan _2 15:04:05 MST 2006” ref time
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if len(rooms) == 0 {
		//	no availability
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
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

	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			Ok:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		res := jsonResponse{
			Ok:      false,
			Message: "error connecting to database",
		}
		out, _ := json.MarshalIndent(res, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	jResp := jsonResponse{
		Ok:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(jResp, "", "     ")
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
		//m.App.ErrorLog.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}
	strMap := make(map[string]string)
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	strMap["start_date"] = sd
	strMap["end_date"] = ed

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
	////		let's get the roomId
	//// we name it id in our rout
	//roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	//if err != nil {
	//	m.App.Session.Put(r.Context(), "error", "can't get roomID from URL")
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}
	// he has changed this to test it more easily
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameters")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation out of the session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

	room, err := m.DB.GetRoomByID(roomID)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "no such room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "login.page.tmpl", r, &models.TemplateData{
		Form: forms.New(nil),
	})
}

//PostShowLogin handles logging user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	//	first thing we should do to secure this login page is to
	//	prevent session fixation attack
	_ = m.App.Session.RenewToken(r.Context())
	//	every session that is stored anywhere in our application
	//	 has a certain token associated it with it
	//	 when u do a login or logout it is a good practice to change these token
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	form := forms.New(r.PostForm)
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, "login.page.tmpl", r, &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)

	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials!")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "logged in successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//LogOut logs out user
func (m *Repository) LogOut(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// remember every time we stop our app and restart it again all of our sessions are gone
// because we use cookies
// in production we may not want to use simple cookies to store our session
// u may want something like Redis 'is perfect for storing sessions'

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "admin-dashboard.page.tmpl", r, &models.TemplateData{})
}
