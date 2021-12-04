package models

import (
	"time"
)

//User Our user table
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Room is the room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Restriction is the Restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Reservation is Reservation  model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Processed int
	Room      Room
}

// RoomRestriction  is RoomRestriction  model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	Room          Room
	ReservationID int
	Reservation   Reservation
	RestrictionID int
	Restrictions  Restriction
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

//MailData is structure of our Email
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
