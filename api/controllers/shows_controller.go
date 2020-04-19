package controllers

import (
	"Dante/api/models"
	"Dante/api/responses"
	"context"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"net/http"
)

func (server *Server) GetShows(w http.ResponseWriter, r *http.Request) {

	shows, err := models.Shows().All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, shows)
}

func (server *Server) GetShow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	show, err := models.Shows(qm.Where("id=?", id)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, show)
}

func (server *Server) GetShowProducers(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	show, err := models.Shows(qm.Where("id=?", id)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	showProducers, err := show.UserIDMemberMembers(qm.Where("id_shows=?", id)).All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, showProducers)

}
