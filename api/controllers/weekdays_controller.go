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
	id := vars["id"]

	showWeekday, err := models.ShowWeekdays(qm.Where("id=?", id)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, showWeekday)

}

func (server *Server) GetShowsWeekDays(w http.ResponseWriter, r *http.Request) {

	showWeekDays, err := models.ShowWeekdays().All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, showWeekDays)
}

func (server *Server) GetShowsOnWeekDay(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	weekDayId := vars["id"]

	showsWeekday, err := models.ShowWeekdays(qm.Where("id_week_days=?", weekDayId)).All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var shows []models.Show
	for _, showWeekday := range showsWeekday {
		show, err := models.Shows(qm.Where("id=?", showWeekday.IDShows)).One(context.Background(), server.DB)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		shows = append(shows, *show)
	}

	responses.JSON(w, http.StatusOK, shows)

}

func (server *Server) GetZonesOnWeekDay(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	weekDayId := vars["id"]

	zonesWeekDay, err := models.DayZones(qm.Where("id_week_days=?", weekDayId)).All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var zones []models.Zone
	for _, zoneWeekDay := range zonesWeekDay {
		zone, err := models.Zones(qm.Where("id=?", zoneWeekDay.IDZones)).One(context.Background(), server.DB)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		zones = append(zones, *zone)
	}

	responses.JSON(w, http.StatusOK, zones)

}
