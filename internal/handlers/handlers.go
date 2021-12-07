package handlers

import (
	"encoding/json"
	"fmt"
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// we need this because we require a specific format for our dare insertion
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")
	//we put it in our string map in our template date
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	log.Println(sd)
	log.Println(ed)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	roomID, err := strconv.Atoi(rID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	room, err := m.DB.GetRoomByID(roomID)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "no such room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		Room:      room,
	}

	//postform has all of the url values and associated data
	form := forms.New(r.PostForm)
	form.Required("firstName", "lastName", "email")
	form.MinLength("firstName", 3)
	form.IsEmail("email")
	if !form.Valid() {

		data := make(map[string]interface{})
		data["reservation"] = reservations
		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed
		//http.Error(w, "invalid form", http.StatusSeeOther)
		render.Template(w, "make-reservation.page.tmpl", r, &models.TemplateData{
			Form:   form,
			Data:   data,
			StrMap: stringMap,
		})
		return
	}
	newReservationID, err := m.DB.InsertReservation(reservations)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation to data base")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	//“Mon Jan _2 15:04:05 MST 2006” ref time
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	//	http.Redirect(w, r, "/", http.StatusSeeOther)
	//	return
	//}
	// he has changed this to test it more easily
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameters")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation out of the session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

//AdminNewReservations Shows all new reservations in admin tools
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.NewReservation()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, "admin-reservations-new.page.tmpl", r, &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservation()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, "admin-reservations-all.page.tmpl", r, &models.TemplateData{
		Data: data,
	})
}

//AdminReservationsCalender displays the reservation calender
func (m *Repository) AdminReservationsCalender(w http.ResponseWriter, r *http.Request) {
	// first we assume there is no month or year specified
	//here we specify params by name not by //
	now := time.Now()

	if r.URL.Query().Get("y") != "" {

		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}
	// here we just get the month and year of next month and prev month
	next := now.AddDate(0, 1, 0)
	prev := now.AddDate(0, -1, 0)
	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")
	prevMonth := prev.Format("01")
	prevMonthYear := prev.Format("2006")

	stringMap := make(map[string]string)
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["prev_month"] = prevMonth
	stringMap["prev_month_year"] = prevMonthYear
	data := make(map[string]interface{})
	data["now"] = now

	//set total days of each month
	//first get the first and last day of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_of_month"] = lastOfMonth.Day()

	rooms, err := m.DB.GetAllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["rooms"] = rooms

	for _, room := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		//	let's loop through month
		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}
		//	get all the restrictions for current room
		restrictions, err := m.DB.GetRestrictionsFroRoomByDate(room.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				//	it is a reservation

				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}

			} else {
				//	 it is a block

				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}
		data[fmt.Sprintf("reservation_map_%d", room.ID)] = reservationMap

		data[fmt.Sprintf("block_map_%d", room.ID)] = blockMap
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", room.ID), blockMap)
	}

	render.Template(w, "admin-reservations-calender.page.tmpl", r, &models.TemplateData{
		StrMap: stringMap,
		Data:   data,
		IntMap: intMap,
	})
}

//PostAdminReservationsCalender Handles Post of reservation calender
func (m *Repository) PostAdminReservationsCalender(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	year, err := strconv.Atoi(r.Form.Get("y"))
	month, err := strconv.Atoi(r.Form.Get("m"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// lets checks for blocks
	rooms, err := m.DB.GetAllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	for _, room := range rooms {
		// remove blocks
		//===========================
		//	let's get the previous state of blocks than compare it
		//	we have put it in the session in get part now we retrieve it and loop through it
		//	if we have any date that is not in posted data and restriction id is greater than zero
		//	 then it is a block we need to remove which its room restrictions id
		//	 moreover if we have a date that is in posted data and not in the map we have to add it
		blockMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.ID)).(map[string]int)
		for name, value := range blockMap {
			// name is the date of restrictions
			//value is its id in db
			//	ok will be false if the value is not in the map
			if val, ok := blockMap[name]; ok {
				//	only pay attention to values gt than zero, and that are not in the post form
				//	the rest are just placeholders for days without blocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", room.ID, name)) {
						err = m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
							return
						}
					}
				}
			}
		}

	}
	//	let's add blocks
	// if the checkbox is checked it will be in posted data otherwise it is not
	for name := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomId, _ := strconv.Atoi(exploded[2])
			//	insert a new block
			//date is separated by - not _
			sd, err := time.Parse("2006-01-2", exploded[3])
			if err != nil {
				helpers.ServerError(w, err)
				return
			}
			block := models.RoomRestriction{

				StartDate: sd,
				EndDate:   sd,
				RoomID:    roomId,
			}
			err = m.DB.InsertBlockForRoom(block)
			if err != nil {
				helpers.ServerError(w, err)
				return
			}

		}
	}
	m.App.Session.Put(r.Context(), "flash", "Changes Saved!")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%d&m=%d", year, month), http.StatusSeeOther)
}

//AdminShowReservation shows the reservation in admin tools
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src
	//we get the reservation
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["year"] = year
	stringMap["month"] = month

	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, "admin-reservations-show.page.tmpl", r, &models.TemplateData{
		StrMap: stringMap,
		Data:   data,
		Form:   forms.New(nil),
	})

}

//PostAdminShowReservation shows the reservation in admin tools
func (m *Repository) PostAdminShowReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src
	//we get the reservation
	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//data := make(map[string]interface{})

	form := forms.New(r.PostForm)
	form.Required("firstName", "lastName", "email")
	form.MinLength("firstName", 3)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = res
		m.App.Session.Put(r.Context(), "error", "Had Some Error!")
		render.Template(w, "admin-reservations-show.page.tmpl", r, &models.TemplateData{
			Form:   form,
			Data:   data,
			StrMap: stringMap,
		})
		return
	}
	res.FirstName = r.Form.Get("firstName")
	res.LastName = r.Form.Get("lastName")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "Changes Saved")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

//AdminProcessReservation makes a reservation processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "error in url!")
		http.Redirect(w, r, fmt.Sprintf("/admin/dashboard"), http.StatusSeeOther)
		return
	}
	src := chi.URLParam(r, "src")

	err = m.DB.UpdateProcessedFroReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	m.App.Session.Put(r.Context(), "flash", "Reservation Marked As Processed!")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", year, month), http.StatusSeeOther)

	}

}

//AdminDeleteReservation makes a reservation gone
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "error in url!")
		http.Redirect(w, r, fmt.Sprintf("/admin/dashboard"), http.StatusSeeOther)
		return
	}
	src := chi.URLParam(r, "src")

	err = m.DB.DeleteReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "error", "A Reservation Has Been Deleted!")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}
