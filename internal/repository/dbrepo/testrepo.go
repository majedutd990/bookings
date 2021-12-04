package dbrepo

import (
	"errors"
	"github.com/majedutd990/bookings/internal/models"
	"log"
	"time"
)

func (p *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (p *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if roomID is 0 this should fail
	if res.RoomID == 0 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

//InsertRoomRestriction inserts a room restriction in room restriction table
func (p *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

//SearchAvailabilityByDatesByRoomID returns true if there is an availability otherwise false for roomID
func (p *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {
	layout := "2006-01-02"
	str := "2049-12-31"
	starDate, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
		return false, err
	}

	//this date should fail no mater what

	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
		return false, err
	}
	// we check it higher than critical start date
	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// any date grater than this should return no availability
	if start.After(starDate) {
		return false, nil
	}

	return true, nil

}

//SearchAvailabilityForAllRooms returns a slice of available room for any date ranges
func (p *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	// if the start date is after 2049-12-31, then return empty slice,
	// indicating no rooms are available;
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return rooms, errors.New("some error")
	}

	if start.After(t) {
		return rooms, nil
	}

	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

//GetRoomByID gets a room by id
func (p *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("there is no such room")
	}
	return room, nil
}

func (p *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u = models.User{
		ID:          0,
		FirstName:   "",
		LastName:    "",
		Email:       "",
		Password:    "",
		AccessLevel: "",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	return u, nil
}

func (p *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (p *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	//correct credential
	//we are faking what happens in db
	if email == "me@here.ca" {
		return 1, "", nil
	}
	return 0, "", errors.New("some error")
}

func (p *testDBRepo) AllReservation() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}
func (p *testDBRepo) NewReservation() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (p *testDBRepo) GetReservationById(id int) (models.Reservation, error) {
	var res models.Reservation
	return res, nil
}

func (p *testDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

func (p *testDBRepo) DeleteReservationById(id int) error {
	return nil
}

func (p *testDBRepo) UpdateProcessedFroReservation(id, processed int) error {
	return nil
}

func (p *testDBRepo) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}
func (p *testDBRepo) GetRestrictionsFroRoomByDate(roomId int, startDate, endDate time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	
	return restrictions, nil
}

func (p *testDBRepo) InsertBlockForRoom(r models.RoomRestriction) error {

	return nil
}

func (p *testDBRepo) DeleteBlockByID(id int) error {

	return nil
}
