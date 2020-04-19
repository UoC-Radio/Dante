package controllers

import (
	"Dante/api/models"
	"Dante/api/responses"
	"Dante/api/utils/formaterror"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (server *Server) CreateShow(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	show := models.Show{}
	err = json.Unmarshal(body, &show)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* validate */
	if show.Title == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("required 'Title'"))
		return
	}
	if show.ProducerNickname == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("required 'ProducerNickname'"))
		return
	}

	err = show.Insert(context.Background(), server.DB, boil.Infer())
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, show)

}

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

func (server *Server) UpdateShow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	/* if input 'id' is valid*/
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	/* if input 'id' exists*/
	exists, err := models.ShowExists(context.Background(), server.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if !exists {
		responses.ERROR(w, http.StatusNotFound, errors.New("id not found"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	showUpdate := models.Show{}
	err = json.Unmarshal(body, &showUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* validate */
	if showUpdate.Title == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("required 'Title'"))
		return
	}
	if showUpdate.ProducerNickname == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("required 'ProducerNickname'"))
		return
	}
	showUpdate.ID = id

	_, err = showUpdate.Update(context.Background(), server.DB, boil.Infer())
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, showUpdate)

}

func (server *Server) DeleteShow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	/* if input 'id' is valid*/
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	/* if input 'id' exists*/
	show, err := models.Shows(qm.Where("id=?", id)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = show.Delete(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", id))
	responses.JSON(w, http.StatusNoContent, "")
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
