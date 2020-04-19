package controllers

import "Dante/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Users routes
	s.Router.HandleFunc("/members/{username}", middlewares.SetMiddlewareIPFilter(s.GetMember)).Methods("GET")

	s.Router.HandleFunc("/weekdays/{id}", middlewares.SetMiddlewareJSON(s.GetShowsWeekDay)).Methods("GET")
	s.Router.HandleFunc("/weekdays/{id}/shows", middlewares.SetMiddlewareJSON(s.GetShowsOnWeekDay)).Methods("GET")
	s.Router.HandleFunc("/weekdays/{id}/zones", middlewares.SetMiddlewareJSON(s.GetZonesOnWeekDay)).Methods("GET")

	s.Router.HandleFunc("/shows", middlewares.SetMiddlewareIPFilter(s.GetShows)).Methods("GET")
	s.Router.HandleFunc("/shows/{id}", middlewares.SetMiddlewareIPFilter(s.GetShow)).Methods("GET")
	s.Router.HandleFunc("/shows/{id}/producers", middlewares.SetMiddlewareIPFilter(s.GetShowProducers)).Methods("GET")

	// Public

}
