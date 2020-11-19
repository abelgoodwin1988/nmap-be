package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-playground/validator/v10"
)

func NewDB(a, p int) (*sql.DB, error) {
	validate := validator.New()
	errs := validate.Var(a, "ip")
	if errs != nil {
		return nil, errs
	}
	errs = validate.Var(p, "hostname_port")
	if errs != nil {
		return nil, errs
	}

	db, err := sql.Open("mysql")
	if err != nil {
		return nil, err
	}
	return db, nil
}
