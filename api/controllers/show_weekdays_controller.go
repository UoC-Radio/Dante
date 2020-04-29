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
	"strings"
	"time"
)

/*
Same POST:
	curl --header "Content-Type: application/json" --request POST --data '{"id_shows":420,"duration":"2h", "start_time":"0000-01-01T12:51:00+00:02"}' http://localhost:8080/weekdays/5/shows
*/
func (server *Server) AddShowOnWeekday(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	weekDayId, err := strconv.Atoi(vars["id"])

	// check that weekday exist
	exists, err := models.WeekDayExists(context.Background(), server.DB, weekDayId)
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

	// schedule check
	showsOnWeekday, _ := models.ShowWeekdays(qm.Where("id_week_days=?", showWeekday.IDWeekDays)).All(context.Background(), server.DB)
	err = checkOverlap(showWeekday, showsOnWeekday)
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
	weekDayId, err := strconv.Atoi(vars["id"])
	showId, err := strconv.Atoi(vars["show_id"])

	// get weekday
	showWeekday, err := models.ShowWeekdays(qm.Where("id_week_days=?", weekDayId), qm.Where("id_shows=?", showId)).One(context.Background(), server.DB)
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

/* helper function */

// checks whether a showWeekDay overlaps in schedule (returns error)
func checkOverlap(showWeekday models.ShowWeekday, showsOnWeekday models.ShowWeekdaySlice) error {

	startTime := showWeekday.StartTime
	duration, _ := time.ParseDuration(showWeekday.Duration)
	endTime := startTime.Add(duration)

	// iterate through all other (scheduled) show_weekdays
	for _, sw := range showsOnWeekday {

		swStart := sw.StartTime
		durationSplit := strings.Split(sw.Duration, ":") //convert sql interval (e.g. "02:00:00") -> golang time.Duration ("2h0m0s")
		swDuration, err := time.ParseDuration(durationSplit[0] + "h" + durationSplit[1] + "m" + durationSplit[2] + "s")
		if err != nil {
			return err
		}
		swEnd := swStart.Add(swDuration)

		if startTime.Before(swEnd) && swStart.Before(endTime) {
			return errors.New("overlaps with show: " + strconv.Itoa(sw.IDShows.Int))
		}

	}

	return nil

}
