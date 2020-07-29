package main

import (
	"log"
	"os"
	"encoding/json"
	"sportsbetting-data-api-chi_router/app"
	"sportsbetting-data-api-chi_router/db"

	"github.com/go-chi/chi"
)

var (
	// EmailConf is the universal email configuration
	EmailConf app.EmailConfig
	_         = json.Unmarshal([]byte(os.Getenv("SBD_API_EMAIL_CONFIG")), &EmailConf)
)

func main() {
	r := chi.NewRouter()
	dbc := &db.Container{
		Conf: db.Dbconfig{
			Rdbms:    "postgres",
			Host:     "local_postgres",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			Database: "data0",
			WinAuth:  false,
		},
	}
	err := dbc.Open()
	if err != nil {
		log.Println(err.Error())
	}
	em := &app.EmailSender{
		Conf: EmailConf,
	}
	server := &app.Server{
		Dbc:         dbc,
		Router:      r,
		EmailSender: em,
	}
	server.Start()
	server.EmailSender.CreateSendEmail()
}
