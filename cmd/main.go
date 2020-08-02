package main

import (
	"log"
	"os"
	"sports-data-api/app"
	"sports-data-api/db"

	"github.com/go-chi/chi"
)

var (
	Rdbms = os.Getenv("SDA_RDBMS")
	Host = os.Getenv("SDA_DB_HOST")
	Port = os.Getenv("SDA_DB_PORT")
	User = os.Getenv("SDA_DB_USER")
	Password = os.Getenv("SDA_DB_PASSWORD")
	Database = os.Getenv("SDA_DATABASE")
)

func main() {
	r := chi.NewRouter()
	dbc := &db.Container{
		Conf: db.Dbconfig{
			Rdbms:    Rdbms,
			Host:     Host,
			Port:     Port,
			User:     User,
			Password: Password,
			Database: Database,
			WinAuth:  false,
		},
	}
	err := dbc.Open()
	if err != nil {
		log.Println(err.Error())
	}
	server := &app.Server{
		Dbc:         dbc,
		Router:      r,
	}
	server.Start()
}
