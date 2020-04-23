package controllers

import (
	"Dante/api/models"
	"Dante/api/responses"
	"context"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"net/http"
)

func (server *Server) GetShowsWeekDay(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	weekDayId := vars["id"]

	showsWeekday, err := models.ShowWeekdays(qm.Where("id_week_days=?", weekDayId)).All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, showsWeekday)

}

func (server *Server) GetZonesWeekDay(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	weekDayId := vars["id"]

	zonesWeekDay, err := models.DayZones(qm.Where("id_week_days=?", weekDayId)).All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, zonesWeekDay)

}
