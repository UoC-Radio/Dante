package controllers

import "Dante/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Users routes
	s.Router.HandleFunc("/members/{username}", middlewares.SetMiddlewareIPFilter(s.GetMember)).Methods("GET")

	// Public

}
