package controllers

import "Dante/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Users routes
	s.Router.HandleFunc("/members/{username}", middlewares.SetMiddlewareIPFilter(s.GetMember)).Methods("GET")

	s.Router.HandleFunc("/weekdays/{id}/shows", middlewares.SetMiddlewareJSON(s.GetShowsWeekDay)).Methods("GET")
	s.Router.HandleFunc("/weekdays/{id}/zones", middlewares.SetMiddlewareJSON(s.GetZonesWeekDay)).Methods("GET")

	s.Router.HandleFunc("/shows", middlewares.SetMiddlewareIPFilter(s.CreateShow)).Methods("POST")
	s.Router.HandleFunc("/shows", middlewares.SetMiddlewareIPFilter(s.GetShows)).Methods("GET")
	s.Router.HandleFunc("/shows/{id}", middlewares.SetMiddlewareIPFilter(s.GetShow)).Methods("GET")
	s.Router.HandleFunc("/shows/{id}", middlewares.SetMiddlewareIPFilter(s.UpdateShow)).Methods("PUT")
	s.Router.HandleFunc("/shows/{id}", middlewares.SetMiddlewareIPFilter(s.DeleteShow)).Methods("DELETE")
	s.Router.HandleFunc("/shows/{id}/producers", middlewares.SetMiddlewareIPFilter(s.GetShowProducers)).Methods("GET")
	s.Router.HandleFunc("/shows/{id}/golive", middlewares.SetMiddlewareIPFilter(s.UpdateGoLive)).Methods("PUT")
	s.Router.HandleFunc("/shows/{id}/activate", middlewares.SetMiddlewareIPFilter(s.SetActiveShow)).Methods("PUT")
	s.Router.HandleFunc("/shows/{id}/deactivate", middlewares.SetMiddlewareIPFilter(s.SetActiveShow)).Methods("PUT")

	s.Router.HandleFunc("/shows/{id}/messages", middlewares.SetMiddlewareIPFilter(s.GetMessages)).Methods("GET")
	s.Router.HandleFunc("/shows/{id}/message", middlewares.SetMiddlewareIPFilter(s.SendMessage)).Methods("POST")

	// Public

}
