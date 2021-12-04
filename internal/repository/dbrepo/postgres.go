package dbrepo

import (
	"context"
	"errors"
	"github.com/majedutd990/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
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

func (p *postgresDBRepo) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	query := `
			select r.id, r.room_name,r.created_at,r.updated_at
			from rooms r
			order by r.room_name
`
	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	var room models.Room
	for rows.Next() {
		err = rows.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
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

//GetUserByID returns user by id
func (p *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	query := `
			select id,first_name,last_name,email,password,access_level,created_at,updated_at from users
			where id = $1
`
	//at most one row
	row := p.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil

}

//UpdateUser updates a user in database
func (p *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	query := `
			update users set first_name=$1,last_name=$2, email=$3, access_level=$4,updated_at=$5
			where id = $6

`
	_, err := p.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil

}

//Authenticate authenticates a user
func (p *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var id int
	var hashedPassword string
	query := `
			select id, password from users where email = $1
`
	row := p.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&id,
		&hashedPassword,
	)
	if err != nil {
		return id, "", err
	}
	// we use a builtin package called bCrypt. We have it in std library and use compare hash and password function
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return id, "", errors.New("incorrect password")
	} else if err != nil {
		return id, "", err
	}
	return id, hashedPassword, nil
}

//AllReservation returns a slice of all reservations
func (p *postgresDBRepo) AllReservation() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var reservations []models.Reservation

	query := `
				select r.id,r.first_name,r.last_name,r.email,r.phone,
				r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at,r.processed
				,rm.id,rm.room_name
				from reservations r 
				left join rooms rm on (r.room_id = rm.id)
				order by r.start_date asc
`
	rows, err := p.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if err := rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

//NewReservation returns a slice of all reservations
func (p *postgresDBRepo) NewReservation() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var reservations []models.Reservation

	query := `
				select r.id,r.first_name,r.last_name,r.email,r.phone,
				r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at,r.processed
				,rm.id,rm.room_name
				from reservations r 
				left join rooms rm on (r.room_id = rm.id)
				where r.processed = 0
				order by r.start_date asc
`
	rows, err := p.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if err := rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

//GetReservationById return one reservation by id
func (p *postgresDBRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var res models.Reservation
	query := `
			select r.id,r.first_name,r.last_name,r.email,r.phone,
		    r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at,r.processed
			,rm.id,rm.room_name
			from reservations r 	
			left join rooms rm on (r.room_id = rm.id)
			where r.id = $1
`
	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

//UpdateReservation updates a reservations in database
func (p *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	query := `
			update reservations set first_name=$1,last_name=$2, email=$3, phone=$4,updated_at=$5
			where id = $6

`
	_, err := p.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil

}

//DeleteReservationById delete one reservations by id
func (p *postgresDBRepo) DeleteReservationById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	query := `
			delete from reservations 
			where id = $1
`
	_, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

//UpdateProcessedFroReservation updates the processed for a reservation by id
func (p *postgresDBRepo) UpdateProcessedFroReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	query := `
			update reservations set processed = $1 
			where id = $2
`
	_, err := p.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}
	return nil
}

//GetRestrictionsFroRoomByDate gets all restrictions of a given room in a given duration
func (p *postgresDBRepo) GetRestrictionsFroRoomByDate(roomId int, startDate, endDate time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var restrictions []models.RoomRestriction
	// coalesce is some kind of shor if
	query := `
				select id,coalesce (reservation_id,0), restriction_id,room_id,start_date,end_date
				from room_restrictions 
				where room_id = $1 and  $2<end_date and $3>= start_date
`

	rows, err := p.DB.QueryContext(ctx, query, roomId, startDate, endDate)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var restriction models.RoomRestriction
		err = rows.Scan(
			&restriction.ID,
			&restriction.ReservationID,
			&restriction.RestrictionID,
			&restriction.RoomID,
			&restriction.StartDate,
			&restriction.EndDate,
		)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, restriction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return restrictions, nil
}

//InsertBlockForRoom inserts a block restrictions for a particular date which is not a reservation
func (p *postgresDBRepo) InsertBlockForRoom(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	stmt := `insert into room_restrictions (start_date,end_date,room_id,restriction_id ,created_at,updated_at)
			  values($1,$2,$3,$4,$5,$6)  returning id`

	_, err := p.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		2,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//DeleteBlockByID deletes a room restrictions
func (p *postgresDBRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	stmt := `delete from room_restrictions
 			where id =$1`

	_, err := p.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
