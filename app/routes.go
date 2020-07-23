package app

import (
	"log"
	"net/http"
	"time"
	"sports-data-api/db"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Server represents a web server object
type Server struct {
	Dbc    *db.Container
	Router *chi.Mux
}

// Routes
func (s *Server) routes() {
	s.Router.Use(
		middleware.RedirectSlashes,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60 * time.Second),
	)
	s.Router.Route("/user", func(r chi.Router) {
		r.Use(middleware.Throttle(10))
		r.Post("/generateToken", s.GenerateToken())						    // working
	})
	s.Router.Route("/api/v1", func(r chi.Router) {
		r.Use(s.Authenticate)
		r.Route("/mlb", func(r chi.Router) {
			r.Get("/teams", s.GetTeams())                                  // working
			r.Get("/baserunning", s.GetBaserunning())                      // working
			r.Get("/baserunning/{teamabbrev}", s.GetBaserunning())        // working
			r.Get("/pitching", s.GetPitching())                            // working
			r.Get("/pitching/{teamabbrev}", s.GetPitching())              // working
			r.Get("/batting", s.GetBatting())                              // working
			r.Get("/batting/{teamabbrev}", s.GetBatting())                // working
			r.Get("/splits/batting", s.GetBattingSplits())                 // working
			r.Get("/splits/batting/{teamabbrev}", s.GetBattingSplits())   // working
			r.Get("/splits/pitching", s.GetPitchingSplits())               // working
			r.Get("/splits/pitching/{teamabbrev}", s.GetPitchingSplits()) // working
		})
	})
}

// Start initializes routes and starts server
func (s *Server) Start() {
	s.routes()
	log.Fatal(http.ListenAndServe(":8600", s.Router))
}
