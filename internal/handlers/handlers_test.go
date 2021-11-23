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
	"reflect"
	"strings"
	"testing"
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
	//{
	//	name:   "post-search-availability",
	//	url:    "/search-availability",
	//	method: "POST",
	//	params: []postData{
	//		{
	//			key:   "start",
	//			value: "2020-01-01",
	//		}, {
	//			key:   "end",
	//			value: "2020-01-02",
	//		},
	//	},
	//	expectedStatusCode: http.StatusOK,
	//}, {
	//	name:   "post-search-availability-json",
	//	url:    "/search-availability-json",
	//	method: "POST",
	//	params: []postData{
	//		{
	//			key:   "start",
	//			value: "2020-01-01",
	//		}, {
	//			key:   "end",
	//			value: "2020-01-02",
	//		},
	//	},
	//	expectedStatusCode: http.StatusOK,
	//},
	//{
	//	name:   "make-reservation-post",
	//	url:    "/make-reservation",
	//	method: "POST",
	//	params: []postData{
	//		{
	//			key:   "firstName",
	//			value: "Majed",
	//		}, {
	//			key:   "lastName",
	//			value: "Nabavian",
	//		}, {
	//			key:   "email",
	//			value: "majedutd@gmail.com",
	//		}, {
	//			key:   "phone",
	//			value: "555-555-5555",
	//		},
	//	},
	//	expectedStatusCode: http.StatusOK,
	//},
}

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

func TestRepository_Reservation(t *testing.T) {
	res := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	//	rr
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", res)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler return wrong response code: got %d, want %d.", rr.Code, http.StatusOK)
	}
	//	test where reservation is not in session(reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler return wrong response code: got %d, want %d.", rr.Code, http.StatusTemporaryRedirect)
	}
	//	test whit nonexistent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	res.RoomID = 100
	session.Put(ctx, "reservation", res)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler return wrong response code: got %d, want %d.", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "firstName=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "lastName=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Smith@John.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=55-555-55")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	rr
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post reservation handler returns wrong response code full form: got %d, want %d.", rr.Code, http.StatusSeeOther)
	}
	//	test for missing body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	rr
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post reservation handler returns wrong response code for missing body: got %d, want %d.", rr.Code, http.StatusTemporaryRedirect)
	}
	//	test for invalid start_date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "firstName=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "lastName=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Smith@John.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=55-555-55")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	rr
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post reservation handler returns wrong response code for invalid start date: got %d, want %d.", rr.Code, http.StatusTemporaryRedirect)
	}
	//	test for invalid end_date
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "firstName=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "lastName=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Smith@John.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=55-555-55")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	rr
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post reservation handler returns wrong response code for invalid end date: got %d, want %d.", rr.Code, http.StatusTemporaryRedirect)
	}
	//	test for invalid room_id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "firstName=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "lastName=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Smith@John.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=55-555-55")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	rr
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post reservation handler returns wrong response code for invalid room id: got %d, want %d.", rr.Code, http.StatusTemporaryRedirect)
	}
	//	test for invalid data
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	// we just make fName -> j which can't be it needs to be 3 so makes the form invalid
	reqBody = fmt.Sprintf("%s&%s", reqBody, "firstName=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "lastName=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Smith@John.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=55-555-55")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	rr
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post reservation handler returns wrong response code for invalid room id: got %d, want %d.", rr.Code, http.StatusSeeOther)
	}
	//	test for failed insert reservation to db
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "firstName=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "lastName=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Smith@John.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=55-555-55")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=0")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	rr
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post reservation handler failed when trying to failed insertion of reservation: got %d, want %d.", rr.Code, http.StatusTemporaryRedirect)
	}
	//	test for failed insert restrictions to db
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "firstName=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "lastName=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Smith@John.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=55-555-55")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	rr
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post reservation handler failed when trying to failed insertion of resetrictions: got %d, want %d.", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr := httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// second case When DB must fail
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// third case  start date in invalid format
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// forth case  end date in invalid format
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	//	fifth if form is nil
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// sixth case rooms are available
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Post availability rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_AvailabilityJson(t *testing.T) {

	//	first case rooms are not available
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "start_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
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
	//	write the rest yourself
	//	here we should not get ok because we set the full date to 2049-12-31
	//	so if we get true there is problem
	if j.Ok {
		t.Error("Got availability when non was expected in Availability json")
	}
	//	second case
	// room should be available but they are not
	reqBody = "start_date=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "start_date=2040-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.AvailabilityJson)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}
	if !j.Ok {
		t.Error("Got not availability when some was expected in Availability json")
	}
	//   third case
	//	 we set the req body to nil and try it one more time
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.AvailabilityJson)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	if j.Ok || j.Message != "Internal server error" {
		t.Error("Got availability when req body was empty!")
	}
	//	fourth case
	//	we make adb error at 2060-01-01
	reqBody = "start_date=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "start_date=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.AvailabilityJson)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}
	if j.Ok || j.Message != "error connecting to database" {
		t.Error("should have got db error but did not get one in Availability json")
	}
}

func TestRepository_ReservationSummary(t *testing.T) {

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	//if we let sd and ed be empty it just does not show them
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
	//	second case reservation not in session
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_ChooseRoom(t *testing.T) {
	//	first case reservation in session
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, err := http.NewRequest("GET", "/choose-room/1", nil)
	if err != nil {
		t.Error(err)
	}
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	//	second case reservation not in session
	req, err = http.NewRequest("GET", "/choose-room/1", nil)
	if err != nil {
		t.Error(err)
	}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	//	Third case bullshit url
	req, _ = http.NewRequest("GET", "/choose-room/fish", nil)
	//if err != nil {
	//	t.Error(err)
	//}
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/fish"

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_BookRoom(t *testing.T) {
	//	first case DB works
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	//	second case db does not work
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=2040-01-02&id=4", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))

	if err != nil {
		log.Println(err)
	}
	return ctx
}
