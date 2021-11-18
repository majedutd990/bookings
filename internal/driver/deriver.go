package driver

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

// the above imports are underlying drivers

//DB holds the database connection pool. its main purpose is to
// be ready for situations where different DBs(mongo, mysql , ...) we are going to use
type DB struct {
	SQL *sql.DB
}

//dbConn connection to our db
var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDBConn = 5

//maxDBLifetime max lifetime for a DB connection
const maxDBLifetime = 5 * time.Minute

//ConnectSql takes a connection string dsn and returns a DB pool for postgres type which is our connection
// it also sets the setting to our setting
func ConnectSql(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		//just die
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetMaxIdleConns(maxIdleDBConn)
	d.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = d
	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// testDB pings DB
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

//NewDatabase create and checks connection
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
