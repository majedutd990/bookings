package dbrepo

import (
	"database/sql"
	"github.com/majedutd990/bookings/internal/config"
	"github.com/majedutd990/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

//let us create a testDb
type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DataBaseRepo {

	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DataBaseRepo {

	return &testDBRepo{
		App: a,
	}
}
