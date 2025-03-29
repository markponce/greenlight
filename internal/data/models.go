package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotfound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
	Users  UserModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
