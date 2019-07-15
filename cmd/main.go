package main

import (
	"baseball_reference_api/app"
	"baseball_reference_api/db"
	"log"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true
	dbc, err := db.NewContainer(db.Dbconfig)
	if err != nil {
		log.Println(err.Error())
	}
	server := &app.Server{
		Dbc:    dbc,
		Router: router,
	}
	server.Start()
}
