package controllers

import "Dante/api/middlewares"

func (s *Server) initializeRoutes() {

	idRegEx := "[0-9]+"

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Users routes
	s.Router.HandleFunc("/members/{username}", middlewares.SetMiddlewareIPFilter(s.GetMember)).Methods("GET")
	s.Router.HandleFunc("/members/{username}/shows", middlewares.SetMiddlewareIPFilter(s.GetMemberShows)).Methods("GET")

	s.Router.HandleFunc("/weekdays/{id:[1-7]}/shows", middlewares.SetMiddlewareJSON(s.GetShowsWeekDay)).Methods("GET")
	s.Router.HandleFunc("/weekdays/{id:[1-7]}/shows", middlewares.SetMiddlewareJSON(s.AddShowOnWeekday)).Methods("POST")
	s.Router.HandleFunc("/weekdays/{id:[1-7]}/shows/{show_id:"+idRegEx+"}", middlewares.SetMiddlewareJSON(s.RemoveShowFromWeekday)).Methods("DELETE")

	s.Router.HandleFunc("/shows", middlewares.SetMiddlewareIPFilter(s.GetShows)).Methods("GET")
	s.Router.HandleFunc("/shows", middlewares.SetMiddlewareIPFilter(s.CreateShow)).Methods("POST")

	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}", middlewares.SetMiddlewareIPFilter(s.GetShow)).Methods("GET")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}", middlewares.SetMiddlewareIPFilter(s.UpdateShow)).Methods("PUT")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}", middlewares.SetMiddlewareIPFilter(s.DeleteShow)).Methods("DELETE")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/producers", middlewares.SetMiddlewareIPFilter(s.GetShowProducers)).Methods("GET")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/producers/{user_id}", middlewares.SetMiddlewareIPFilter(s.AddOrRemoveShowProducer)).Methods("PUT", "DELETE")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/urls", middlewares.SetMiddlewareIPFilter(s.GetShowUrls)).Methods("GET")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/urls", middlewares.SetMiddlewareIPFilter(s.AddShowUrl)).Methods("POST")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/urls", middlewares.SetMiddlewareIPFilter(s.RemoveShowUrl)).Methods("DELETE")

	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/golive", middlewares.SetMiddlewareIPFilter(s.UpdateGoLive)).Methods("PUT")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/activate", middlewares.SetMiddlewareIPFilter(s.SetActiveShow)).Methods("PUT")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/deactivate", middlewares.SetMiddlewareIPFilter(s.SetActiveShow)).Methods("PUT")

	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/messages", middlewares.SetMiddlewareIPFilter(s.GetMessages)).Methods("GET").Queries("page", "{page:[1-9][0-9]*}").Name("Pagination")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/messages", middlewares.SetMiddlewareIPFilter(s.GetMessages)).Methods("GET")
	s.Router.HandleFunc("/shows/{id:"+idRegEx+"}/messages", middlewares.SetMiddlewareIPFilter(s.SendMessage)).Methods("POST")
}
