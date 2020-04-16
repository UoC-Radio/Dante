package controllers

import (
	"Dante/api/models"
	"Dante/api/responses"
	"context"
	_ "encoding/json"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"net/http"
)

func (server *Server) GetMember(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	username := vars["username"]
	member, err := models.Members(qm.Where("username=?", username)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, member)
}
