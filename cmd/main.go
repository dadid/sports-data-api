package main

import (
	"encoding/json"
	"log"
	"os"
	"sportsbetting-data-api/app"
	"sportsbetting-data-api/db"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

var (
	dbconf db.Dbconfig
)

func main() {
	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true
	err := json.Unmarshal([]byte(os.Getenv("SBD_API_DB_CONFIG")), &dbconf)
	if err != nil {
		log.Println(errors.Wrap(err, "error unsharmshalling database config from environment variable").Error())
	}
	dbc := &db.Container{
		Conf: dbconf,
	}
	err = dbc.Open()
	if err != nil {
		log.Println(err.Error())
	}
	server := &app.Server{
		Dbc:    dbc,
		Router: router,
	}
	server.Start()
}
