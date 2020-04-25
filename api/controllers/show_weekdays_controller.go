package controllers

import (
	"Dante/api/models"
	"Dante/api/responses"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (server *Server) AddShowOnWeekday(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	weekDayId, err := strconv.Atoi(vars["id"])

	// check that weekday exist
	exists, err := models.ShowWeekdayExists(context.Background(), server.DB, weekDayId)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if !exists {
		responses.ERROR(w, http.StatusBadRequest, errors.New("weekday does not exist"))
		return
	}

	// read request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	showWeekday := models.ShowWeekday{}
	err = json.Unmarshal(body, &showWeekday)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	showWeekday.IDWeekDays = null.Int{weekDayId, true}

	//check that show exist
	exists, err = models.ShowExists(context.Background(), server.DB, showWeekday.IDShows.Int)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if !exists {
		responses.ERROR(w, http.StatusNotFound, errors.New("show does not exist"))
		return
	}

	//parse duration
	_, err = time.ParseDuration(showWeekday.Duration)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = showWeekday.Insert(context.Background(), server.DB, boil.Infer())
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, showWeekday)

}

func (server *Server) RemoveShowFromWeekday(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	weekDayId, err := strconv.Atoi(vars["id_weekday"])
	showID, err := strconv.Atoi(vars["id_show"])

	// get weekday
	showWeekday, err := models.ShowWeekdays(qm.Where("id_week_days=?", weekDayId), qm.Where("id_shows=?", showID)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// delete show_weekday
	_, err = showWeekday.Delete(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, "")

}
