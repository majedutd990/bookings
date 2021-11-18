package dbrepo

import (
	"context"
	"github.com/majedutd990/bookings/internal/models"
	"time"
)

func (p *postgresDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (p *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	var newID int
	//if in the middle of transaction something happens to the user's connection
	// and he or she closes the browser while the transaction is still going on
	// we don't want this to happen we will use context instead happened in go 1.8
	// we use cancel to cancel the context if something goes wrong
	// a context that is always available anywhere in our application called context.background
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	//the above line means that if this transaction is not commit in 4 second he says 3
	//something is seriously wrong in our application
	stmt := `insert into reservations (first_name,last_name,email,phone,start_date,end_date,room_id
             ,created_at,updated_at)
			  values($1,$2,$3,$4,$5,$6,$7,$8,$9)  returning id`

	//exec does not know anything about context but
	// execContext know
	err := p.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now()).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

//InsertRoomRestriction inserts a room restriction in room restriction table
func (p *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	stmt := `insert into room_restrictions (start_date,end_date,room_id,
             reservation_id,created_at,updated_at,restriction_id)
			  values($1,$2,$3,$4,$5,$6,$7)`

	//exec does not know anything about context but
	// execContext know
	_, err := p.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID)

	if err != nil {
		return err
	}
	return nil
}

//SearchAvailabilityByDatesByRoomID returns true if there is an availability otherwise false for roomID
func (p *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	query := `
			select
				count(id)
			from room_restrictions
			where
			room_id = $1 and $2<end_date and $3>start_date`
	var nomRows int
	row := p.DB.QueryRowContext(ctx, query, roomId, start, end)
	err := row.Scan(&nomRows)
	if err != nil {
		return false, err
	}
	if nomRows == 0 {
		return true, nil
	}
	return false, nil

}

//SearchAvailabilityForAllRooms returns a slice of available room for any date ranges
func (p *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `
			select r.id, r.room_name
			from
			rooms r
			where r.id not in 
			(select room_id from room_restrictions rr where $1<rr.end_date and $2>rr.start_date)
`
	rows, err := p.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	var room models.Room
	for rows.Next() {
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

//GetRoomByID gets a room by id
func (p *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var room models.Room

	query := `
			select id,room_name,created_at,updated_at from rooms
			where id = $1
`
	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}
	return room, nil
}
