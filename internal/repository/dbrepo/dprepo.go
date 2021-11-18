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

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DataBaseRepo {

	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
