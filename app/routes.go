package app

import (
	"baseball_reference_api/db"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Server represents a web server object
type Server struct {
	Dbc    db.Container
	Router *httprouter.Router
	// email  *EmailSender
}

// Routes
func (s *Server) routes() {
	defaultMiddleware := []Middleware{
		s.addRequestID,
		s.authenticate,
	}
	s.Router.POST("/api/v1/generateToken", s.httpRouterHandleWrapper(s.limitNumClients(s.generateToken(), 10)))                                                          // working
	s.Router.GET("/api/v1/mlb/teams", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetAllTeams), defaultMiddleware...)))                               // working
	s.Router.GET("/api/v1/mlb/baserunning", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetAllBaserunning), defaultMiddleware...)))                   // working
	s.Router.GET("/api/v1/mlb/baserunning/:teamabbrev", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetTeamBaserunning), defaultMiddleware...)))      // working
	s.Router.GET("/api/v1/mlb/pitching", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetAllPitching), defaultMiddleware...)))                         // working
	s.Router.GET("/api/v1/mlb/pitching/:teamabbrev", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetTeamPitching), defaultMiddleware...)))            // working
	s.Router.GET("/api/v1/mlb/batting", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetAllBatting), defaultMiddleware...)))                           // working
	s.Router.GET("/api/v1/mlb/batting/:teamabbrev", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetTeamBatting), defaultMiddleware...)))              // working
	s.Router.GET("/api/v1/mlb/splits/batting", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetAllBattingSplits), defaultMiddleware...)))              // working
	s.Router.GET("/api/v1/mlb/splits/batting/:teamabbrev", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetTeamBattingSplits), defaultMiddleware...))) // working
	s.Router.GET("/api/v1/mlb/splits/pitching", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetAllPitchingSplits), defaultMiddleware...)))            // working
	s.Router.GET("/api/v1/mlb/splits/pitching/:teamabbrev", s.httpRouterHandleWrapper(s.chainMiddleware(http.HandlerFunc(s.GetTeamPitchingSplits), defaultMiddleware...)))
}

// Start initializes routes and starts server
func (s *Server) Start() {
	s.routes()
	log.Fatal(http.ListenAndServe(":8600", s.Router))
}
