package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/majedutd990/bookings/internal/driver"
	"github.com/majedutd990/bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{
		name:               "home",
		url:                "/",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "about",
		url:                "/about",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "generals-quarters",
		url:                "/generals-quarters",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "majors-suites",
		url:                "/majors-suites",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "search-availability",
		url:                "/search-availability",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "contact",
		url:                "/contact",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "non-existent",
		url:                "/green/egs/ham",
		method:             "GET",
		expectedStatusCode: http.StatusNotFound,
	},
	{
		name:               "login",
		url:                "/user/login",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "logout",
		url:                "/user/logout",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "dashboard",
		url:                "/admin/dashboard",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "reservations_new",
		url:                "/admin/reservations-new",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "reservations_all",
		url:                "/admin/reservations-all",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	}, {
		name:               "show_reservations",
		url:                "/admin/reservations/new/1/show",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "show res cal",
		url:                "/admin/reservations-calender",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "show res cal with param",
		url:                "/admin/reservations-calender?y=2021&m=10",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
	},
}

//TestHandlers tests all routes that are only get requests
func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//	create a test server
	ts := httptest.NewTLSServer(routes)
	//	we close this server when we are done with it
	// defer close the ts after the testHandlers is finished
	// defer actually makes the closing parts wait
	defer ts.Close()
	for _, test := range theTests {
		resp, err := ts.Client().Get(ts.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatal()
		}
		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s expected %d but got %d status code!", test.name, test.expectedStatusCode, resp.StatusCode)
		}

	}
}

var reservationTests = []struct {
	name               string
	reservation        models.Reservation
	expectedStatusCode int
	expectedLocation   string
	expectedHtml       string
}{
	{
		name: "reservation in session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       0,
				RoomName: "General's Quarters",
			},
		},
		expectedStatusCode: http.StatusOK,
		expectedLocation:   "",
		expectedHtml:       `action="/make-reservation"`,
	},
	{
		name:               "reservation not in session",
		reservation:        models.Reservation{},
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
		expectedHtml:       "",
	},
	{
		name: "Room Id Not In DB",
		reservation: models.Reservation{
			RoomID: 100,
		},
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
		expectedHtml:       "",
	},
}

//TestRepository_Reservation tests the reservation's handler with url /make-reservation and uses the above table test: reservationTests
func TestRepository_Reservation(t *testing.T) {
	for _, e := range reservationTests {
		req, _ := http.NewRequest("GET", "/make-reservation", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}
		handler := http.HandlerFunc(Repo.Reservation)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatusCode {
			t.Errorf("Reservation Test Failed: expected code: %d got %d", e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {

			actualLocation, _ := rr.Result().Location()
			log.Println(e.name)
			if actualLocation.String() != e.expectedLocation {
				t.Errorf("Reservation Test Failed: expected location: %s went %s", e.expectedLocation, actualLocation.String())
			}
		}

		if e.expectedHtml != "" {
			expHtml := rr.Body.String()

			if !strings.Contains(expHtml, e.expectedHtml) {
				t.Errorf("Reservation Test Failed: expected html: %s got %s", e.expectedHtml, expHtml)
			}
		}

	}
}

var postReservationTests = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHtml         string
}{
	{
		name: "valid data",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"firstName":  {"John"},
			"lastName":   {"Smith"},
			"email":      {"Smith@John.com"},
			"phone":      {"55-555-55"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/reservation-summary",
		expectedHtml:         "",
	},
	{
		name:                 "no body",
		postedData:           nil,
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/",
		expectedHtml:         "",
	},
	{
		name: "invalid start date",
		postedData: url.Values{
			"start_date": {"invalid"},
			"end_date":   {"2050-01-02"},
			"firstName":  {"John"},
			"lastName":   {"Smith"},
			"email":      {"Smith@John.com"},
			"phone":      {"55-555-55"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/",
		expectedHtml:         "",
	}, {
		name: "invalid end date",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"invalid"},
			"firstName":  {"John"},
			"lastName":   {"Smith"},
			"email":      {"Smith@John.com"},
			"phone":      {"55-555-55"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/",
		expectedHtml:         "",
	},
	{
		name: "invalid room_id",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"firstName":  {"John"},
			"lastName":   {"Smith"},
			"email":      {"Smith@John.com"},
			"phone":      {"55-555-55"},
			"room_id":    {"invalid"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/",
		expectedHtml:         "",
	},
	{
		name: "invalid form",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"firstName":  {"J"},
			"lastName":   {"Smith"},
			"email":      {"Smith@John.com"},
			"phone":      {"55-555-55"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusOK,
		expectedLocation:     "",
		expectedHtml:         `action="/make-reservation"`,
	},
	{
		name: "failed insertion reservation",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"firstName":  {"John"},
			"lastName":   {"Smith"},
			"email":      {"Smith@John.com"},
			"phone":      {"55-555-55"},
			//0 is making an error in test db repo
			"room_id": {"0"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/",
		expectedHtml:         "",
	},
	{
		name: "failed to insert room restriction",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"firstName":  {"John"},
			"lastName":   {"Smith"},
			"email":      {"Smith@John.com"},
			"phone":      {"55-555-55"},
			//1000 is making an error in test db repo
			"room_id": {"1000"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/",
		expectedHtml:         "",
	},
}

func TestRepository_PostReservation(t *testing.T) {
	for _, e := range postReservationTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/make-reservation", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("PostReservation Falied test name %s : expected Code %d got %d", e.name, e.expectedResponseCode, rr.Code)
		}

		if e.expectedLocation != "" {
			actualLocation, _ := rr.Result().Location()

			if actualLocation.String() != e.expectedLocation {
				t.Errorf("PostReservation Falied test name %s : expected location %s got %s", e.name, e.expectedLocation, actualLocation.String())
			}
		}

		if e.expectedHtml != "" {
			actualHtml := rr.Body.String()

			if !strings.Contains(actualHtml, e.expectedHtml) {
				t.Errorf("PostReservation Falied test name %s : expected html %s got %s", e.name, e.expectedHtml, actualHtml)
			}
		}
	}
}

var testAvailabilityData = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name: "rooms not available",
		postedData: url.Values{
			"start": {"2050-01-01"},
			"end":   {"2050-01-02"},
		},
		expectedResponseCode: http.StatusSeeOther,
	}, {
		name: "rooms are available",
		postedData: url.Values{
			"start":   {"2040-01-01"},
			"end":     {"2040-01-02"},
			"room_id": {"1"},
		},
		expectedResponseCode: http.StatusOK,
		expectedLocation:     "",
	},
	{
		name: "DB Must fail",
		postedData: url.Values{
			"start": {"2060-01-01"},
			"end":   {"2060-01-02"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/",
	},
	{
		name: "invalid start date",
		postedData: url.Values{
			"start": {"invalid"},
			"end":   {"2040-01-02"},
		},
		expectedResponseCode: http.StatusSeeOther,
	}, {
		name: "invalid end date",
		postedData: url.Values{
			"start": {"2040-01-01"},
			"end":   {"invalid"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
	{
		name:                 "form is nil",
		postedData:           nil,
		expectedResponseCode: http.StatusSeeOther,
	},
}

func TestRepository_PostAvailability(t *testing.T) {
	for _, e := range testAvailabilityData {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/search-availability", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostAvailability)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("Search Availability Falied test name %s : expected Code %d got %d", e.name, e.expectedResponseCode, rr.Code)
		}

		if e.expectedLocation != "" {
			actualLocation, _ := rr.Result().Location()

			if actualLocation.String() != e.expectedLocation {
				t.Errorf("Search Availability Falied test name %s : expected location %s got %s", e.name, e.expectedLocation, actualLocation.String())
			}
		}
	}
}

var testAvailabilityJsonData = []struct {
	name            string
	postedData      url.Values
	expectedOk      bool
	expectedMessage string
}{
	{
		name: "rooms not available",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-01"},
			"room_id":    {"1"},
		},
		expectedOk:      false,
		expectedMessage: "",
	}, {
		name: "room should be available but they are not",
		postedData: url.Values{
			"start_date": {"2040-01-01"},
			"end_date":   {"2040-01-01"},
			"room_id":    {"1"},
		},
		expectedOk:      true,
		expectedMessage: "",
	}, {
		name:            "req body nil",
		postedData:      nil,
		expectedOk:      false,
		expectedMessage: "Internal server error",
	}, {
		name: "DB Error",
		postedData: url.Values{
			"start_date": {"2060-01-01"},
			"end_date":   {"2060-01-01"},
			"room_id":    {"1"},
		},
		expectedOk:      false,
		expectedMessage: "error connecting to database",
	},
}

func TestRepository_AvailabilityJson(t *testing.T) {

	for _, e := range testAvailabilityJsonData {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/search-availability-json", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler := http.HandlerFunc(Repo.AvailabilityJson)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		var j jsonResponse
		err := json.Unmarshal([]byte(rr.Body.String()), &j)
		if err != nil {
			t.Error("Failed to parse json:AvailabilityJson")
		}
		if j.Ok != e.expectedOk {
			t.Errorf(" Availability json Name %s : expected Response %t  Got  response %t", e.name, e.expectedOk, j.Ok)
		}

		if e.expectedMessage != "" {
			if e.expectedMessage != j.Message {
				t.Errorf(" Availability json Name %s : expected Message %s  Got  message %s", e.name, e.expectedMessage, j.Message)
			}
		}

	}
}

var redservationSummaryTest = []struct {
	name               string
	reservation        models.Reservation
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name: "res in session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "General's Quarters",
			},
		},
		url:                "/reservation-summary",
		expectedStatusCode: http.StatusOK,
		expectedLocation:   "",
	},
	{
		name:               "res not in session",
		reservation:        models.Reservation{},
		url:                "/reservation-summary",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
}

func TestRepository_ReservationSummary(t *testing.T) {

	for _, e := range redservationSummaryTest {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.ReservationSummary)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatusCode {
			t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLocation, _ := rr.Result().Location()
			if actualLocation.String() != e.expectedLocation {
				t.Errorf("ReservationSummary handler returned wrong location : got %s, wanted %s", actualLocation.String(), e.expectedLocation)
			}
		}
	}
}

var chooseRoomTests = []struct {
	name               string
	reservation        models.Reservation
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name: "reservation-in-session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "General's Quarters",
			},
		},
		url:                "/choose-room/1",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/make-reservation",
	},
	{
		name:               "reservation-not-in-session",
		reservation:        models.Reservation{},
		url:                "/choose-room/1",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name:               "malformed-url",
		reservation:        models.Reservation{},
		url:                "/choose-room/fish",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
}

func TestRepository_ChooseRoom(t *testing.T) {
	for _, e := range chooseRoomTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		// set the RequestURI on the request so that we can grab the ID from the URL
		req.RequestURI = e.url

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.ChooseRoom)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var bookRoomTests = []struct {
	name               string
	url                string
	expectedStatusCode int
}{
	{
		name:               "database-works",
		url:                "/book-room?s=2050-01-01&e=2050-01-02&id=1",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "database-fails",
		url:                "/book-room?s=2040-01-01&e=2040-01-02&id=4",
		expectedStatusCode: http.StatusSeeOther,
	},
}

func TestRepository_BookRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	for _, e := range bookRoomTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		session.Put(ctx, "reservation", reservation)

		handler := http.HandlerFunc(Repo.BookRoom)

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("%s failed: returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

//create some data for login tests
// we don't butter our self with password
//expectedStatusCode is clear
//expectedHtml       string pull html out of response and look for things that should be there
//expectedLocation   string what url does the user have
var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHtml       string
	expectedLocation   string
}{
	{
		name:               "valid-credentials",
		email:              "me@here.ca",
		expectedStatusCode: http.StatusSeeOther,
		expectedHtml:       "",
		expectedLocation:   "/",
	}, {
		//go check test_repo
		name:               "Invalid-Credential",
		email:              "jack@nimble.com",
		expectedStatusCode: http.StatusSeeOther,
		expectedHtml:       "",
		expectedLocation:   "/user/login",
	}, {
		name:  "Invalid-Data",
		email: "jj",
		//because we are not doing a redirect we are doing a render, so it is status ok
		expectedStatusCode: http.StatusOK,
		expectedHtml:       `action="/user/login"`,
		expectedLocation:   "",
	},
}

func TestRepository_ShowLogin(t *testing.T) {
	for _, e := range loginTests {
		//	let's create some posted data
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "123hj123")
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		//	 call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf(fmt.Sprintf("Show Login:%s failed expected %d Got %d", e.name, e.expectedStatusCode, rr.Code))
		}
		//	what about expected location

		if e.expectedLocation != "" {
			//	we get the returned location
			location, _ := rr.Result().Location()

			if location.String() != e.expectedLocation {
				t.Errorf("Show Login Failed: %s expected location is:%s went to %s.", e.name, e.expectedLocation, location.String())
			}
		}
		//checking for expected html
		if e.expectedHtml != "" {
			//	we get the returned html
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHtml) {
				t.Errorf("Show Login Failed: %s expected html is:%s got %s.", e.name, e.expectedHtml, html)
			}
		}

	}

}

var adminPostReservationTest = []struct {
	name                 string
	url                  string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name:                 "form is nil",
		url:                  "/admin/reservations/new/1/show",
		postedData:           nil,
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/dashboard",
		expectedHTML:         "",
	}, {
		name: "valid data from new",
		url:  "/admin/reservations/new/1/show",
		postedData: url.Values{
			"firstName": {"John"},
			"lastName":  {"Smith"},
			"email":     {"Smith@John.com"},
			"phone":     {"55-555-55"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-new",
		expectedHTML:         "",
	}, {
		name: "valid data from all",
		url:  "/admin/reservations/all/1/show",
		postedData: url.Values{
			"firstName": {"John"},
			"lastName":  {"Smith"},
			"email":     {"Smith@John.com"},
			"phone":     {"55-555-55"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-all",
		expectedHTML:         "",
	},
	{
		name: "valid data from cal",
		url:  "/admin/reservations/cal/1/show",
		postedData: url.Values{
			"firstName": {"John"},
			"lastName":  {"Smith"},
			"email":     {"Smith@John.com"},
			"phone":     {"55-555-55"},
			"year":      {"2022"},
			"month":     {"01"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-calender?y=2022&m=01",
		expectedHTML:         "",
	},
}

func TestRepository_PostAdminShowReservation(t *testing.T) {
	for _, e := range adminPostReservationTest {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", e.url, strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", e.url, nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.RequestURI = e.url

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostAdminShowReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var adminPostReservationCalendarTests = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
	blocks               int
	reservations         int
}{
	{
		name: "cal",
		postedData: url.Values{
			"y": {time.Now().Format("2006")},
			"m": {time.Now().Format("01")},
			fmt.Sprintf("add_block_1_%s", time.Now().Format("2006-01-2")): {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     fmt.Sprintf("/admin/reservations-calender?y=%s&m=%s", time.Now().Format("2006"), time.Now().Format("01")),
		expectedHTML:         "",
		blocks:               0,
		reservations:         0,
	},
	{
		name: "cal-blocks",
		postedData: url.Values{
			"y": {time.Now().Format("2006")},
			"m": {time.Now().Format("01")},
		},
		expectedResponseCode: http.StatusSeeOther,
		blocks:               1,
	},
	{
		name: "cal-res",
		postedData: url.Values{
			"y": {time.Now().Format("2006")},
			"m": {time.Now().Format("01")},
		},
		expectedResponseCode: http.StatusSeeOther,
		reservations:         1,
	},
}

func TestRepository_PostAdminReservationsCalender(t *testing.T) {

	for _, e := range adminPostReservationCalendarTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/reservations-calender", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/reservations-calender", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		now := time.Now()
		//block mp from date to its restriction id
		bm := make(map[string]int)
		//reservation map
		rm := make(map[string]int)

		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			rm[d.Format("2006-01-2")] = 0
			bm[d.Format("2006-01-2")] = 0
		}

		// if test is a block
		if e.blocks == 1000 {
			bm[now.Format("2006-01-2")] = e.blocks
		}
		if e.blocks > 0 {
			bm[firstOfMonth.Format("2006-01-2")] = e.blocks
		}
		// if test is a reservation
		if e.reservations > 0 {
			rm[lastOfMonth.Format("2006-01-2")] = e.reservations
		}
		session.Put(ctx, "block_map_1", bm)
		session.Put(ctx, "reservation_map_1", rm)
		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler := http.HandlerFunc(Repo.PostAdminReservationsCalender)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed postAdminCalebder %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

var adminDeleteReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "delete-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "delete-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminDeleteReservation(t *testing.T) {
	for _, e := range adminDeleteReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process/delete/cal/1/do%s", e.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

var adminProcessReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "process-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "process-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestRepository_AdminProcessReservation(t *testing.T) {

	for _, e := range adminProcessReservationTests {

		var req *http.Request
		req, _ = http.NewRequest("GET", fmt.Sprintf("/admin/process/reservation/cal/1/do%s", e.queryParams), nil)
		log.Println(req.URL.RequestURI())
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminProcessReservation)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))

	if err != nil {
		log.Println(err)
	}
	return ctx
}
