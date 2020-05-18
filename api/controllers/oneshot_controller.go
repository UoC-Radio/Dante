package controllers

import (
	"Dante/api/models"
	"Dante/api/responses"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"io/ioutil"
	"net/http"
	"time"
)

func (server *Server) CreateOneShot(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	showOneShot := models.ShowOneshot{}
	err = json.Unmarshal(body, &showOneShot)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// check that show exists
	exists, err := models.ShowExists(context.Background(), server.DB, showOneShot.IDShows.Int)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if !exists {
		responses.ERROR(w, http.StatusNotFound, errors.New("show id not found"))
		return
	}

	// check scheduled time
	if showOneShot.ScheduledTime.IsZero() { // TODO: throw error here ?
		showOneShot.ScheduledTime = time.Now()
	}

	// check duration
	if showOneShot.Duration == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("required field 'duration'"))
		return
	}
	duration, err := time.ParseDuration(showOneShot.Duration)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	if duration == 0 {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("'duration' is zero"))
		return
	}

	err = showOneShot.Insert(context.Background(), server.DB, boil.Infer())
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, showOneShot)

}

func (server *Server) GetOneShot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	oneShot, err := models.ShowOneshots(qm.Where("id=?", id)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, oneShot)
}

func (server *Server) GetAllOneShots(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (server *Server) DeleteOneShot(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	oneShotId := vars["id"]

	/* if input 'id' exists*/
	oneShot, err := models.ShowOneshots(qm.Where("id=?", oneShotId)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = oneShot.Delete(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, "")

}

func (server *Server) GetShowOneShots(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	showId := vars["id"]

	show, err := models.Shows(qm.Where("id=?", showId)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	showOneShots, err := show.IDShowShowOneshots(qm.Where("id_shows=?", showId)).All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, showOneShots)

}
