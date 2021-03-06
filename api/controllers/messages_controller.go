package controllers

import (
	"Dante/api/models"
	"Dante/api/responses"
	"Dante/api/utils/formaterror"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

/*
Sample post:
curl --header "Content-Type: application/json" \
--request POST --data '{"nickname":"mitsos","message":"Hello, World!"}' \
http://localhost:8080/shows/444/message
*/
func (server *Server) SendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	showId, err := strconv.Atoi(vars["id"])

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	message := models.ShowMessage{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Fill extra fields
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	message.IPAddr = ip
	message.UserAgent = r.UserAgent()
	message.IDShows = null.IntFrom(showId)

	err = message.Insert(context.Background(), server.DB, boil.Infer())

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, message.ID))
	responses.JSON(w, http.StatusCreated, message)
}

func (server *Server) GetMessages(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	show_id := vars["id"]
	limit, found := vars["limit"]

	var messages models.ShowMessageSlice
	var err error

	if found {
		limit, err := strconv.Atoi(limit)

		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		messages, err = models.ShowMessages(qm.Where("id_shows=?", show_id), qm.OrderBy("received_datetime DESC"), qm.Limit(limit), qm.Offset(20)).All(context.Background(), server.DB)
	} else {
		messages, err = models.ShowMessages(qm.Where("id_shows=?", show_id), qm.OrderBy("received_datetime DESC")).All(context.Background(), server.DB)
	}

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, messages)
}
