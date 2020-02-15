package main

import (
	"log"
	"sportsbetting-data-api-chi_router/app"
	"sportsbetting-data-api-chi_router/db"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	dbc := &db.Container{
		Conf: db.Dbconfig{
			Rdbms:    "postgres",
			Host:     "local_postgres",
			Port:     "5432",
			User:     "postgres",
			Password: "Hotdog10!",
			Database: "data0",
			WinAuth:  false,
		},
	}
	err := dbc.Open()
	if err != nil {
		log.Println(err.Error())
	}
	server := &app.Server{
		Dbc:    dbc,
		Router: r,
	}
	server.Start()
}
