package controllers

import (
	"net/http"

	"Dante/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Root entry for the API")
}
