package repository

import (
	"github.com/majedutd990/bookings/internal/models"
	"time"
)

type DataBaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, rID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetAllRooms() ([]models.Room, error)

	//users function

	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	//Admin function

	AllReservation() ([]models.Reservation, error)
	NewReservation() ([]models.Reservation, error)
	GetReservationById(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservationById(id int) error
	UpdateProcessedFroReservation(id, processed int) error
	GetRestrictionsFroRoomByDate(roomId int, startDate, endDate time.Time) ([]models.RoomRestriction, error)
	InsertBlockForRoom(r models.RoomRestriction) error
	DeleteBlockByID(id int) error
}
