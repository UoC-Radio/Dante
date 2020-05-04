package controllers

import (
	"Dante/api/models"
	"Dante/api/responses"
	"Dante/api/utils/formaterror"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
Sample POST :
curl -X POST -H 'Content-Type: multipart/form-data' -F "title=testTitle" -F "producerNickname=testNick" -F "logoFile=@image.png" localhost:8080/shows
*/
func (server *Server) CreateShow(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	show := models.Show{}

	/* parse multipart-form values */
	decoder := schema.NewDecoder()
	err = decoder.Decode(&show, r.PostForm)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if show.Title == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("required 'Title'"))
		return
	}

	if show.ProducerNickname == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("required field 'ProducerNickname'"))
		return
	}

	// if uq 'title' already exists
	n, err := models.Shows(qm.Where("title=?", show.Title)).Count(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if n > 0 {
		responses.ERROR(w, http.StatusBadRequest, errors.New("'title' field value already exists"))
		return
	}

	/* handle logo filename (multipart-form file) */
	logoFile, handler, err := r.FormFile("logoFile")
	logoFileName := ""
	if logoFile != nil {
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		/* store logo image */
		logoFileName, err = storeLogoImage(logoFile, handler)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}

	}
	show.LogoFilename = null.String{String: logoFileName, Valid: true}

	err = show.Insert(context.Background(), server.DB, boil.Infer())
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusCreated, show)

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

func (server *Server) AddOrRemoveShowProducer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	showId := vars["id"]
	memberId := vars["user_id"]

	show, err := models.Shows(qm.Where("id=?", showId)).One(context.Background(), server.DB)
	if err != nil {
		if show == nil {
			responses.ERROR(w, http.StatusBadRequest, errors.New("show id not found"))
		} else {
			responses.ERROR(w, http.StatusBadRequest, err)
		}
		return
	}

	member, err := models.Members(qm.Where("user_id=?", memberId)).One(context.Background(), server.DB)
	if err != nil {
		if member == nil {
			responses.ERROR(w, http.StatusBadRequest, errors.New("member id not found"))
		} else {
			responses.ERROR(w, http.StatusBadRequest, err)
		}
		return
	}

	// PUT or DELETE show_producer
	returnStatus := http.StatusBadRequest
	if r.Method == "PUT" {
		err = show.AddUserIDMemberMembers(context.Background(), server.DB, false, member)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		returnStatus = http.StatusOK
	} else if r.Method == "DELETE" {
		err = show.RemoveUserIDMemberMembers(context.Background(), server.DB, member)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		returnStatus = http.StatusNoContent
	}
	responses.JSON(w, returnStatus, "")
}

func (server *Server) AddShowUrl(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	showId, err := strconv.Atoi(vars["id"])

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	show, err := models.Shows(qm.Where("id=?", showId)).One(context.Background(), server.DB)
	if err != nil {
		if show == nil {
			responses.ERROR(w, http.StatusBadRequest, errors.New("show id not found"))
		} else {
			responses.ERROR(w, http.StatusBadRequest, err)
		}
		return
	}

	showUrl := models.ShowURL{}
	err = json.Unmarshal(body, &showUrl)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if showUrl.URLURI == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("required 'url_uri'"))
		return
	}
	showUrl.IDShows = null.Int{showId, true}

	err = showUrl.Insert(context.Background(), server.DB, boil.Infer())
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, showUrl)

}

func (server *Server) GetShowUrls(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	showId, err := strconv.Atoi(vars["id"])

	show, err := models.Shows(qm.Where("id=?", showId)).One(context.Background(), server.DB)
	if err != nil {
		if show == nil {
			responses.ERROR(w, http.StatusBadRequest, errors.New("show id not found"))
		} else {
			responses.ERROR(w, http.StatusBadRequest, err)
		}
		return
	}

	showUrls, err := show.IDShowShowUrls(qm.Where("id_shows=?", showId)).All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, showUrls)

}

func (server *Server) RemoveShowUrl(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	showId, err := strconv.Atoi(vars["id"])

	showUrl, err := models.ShowUrls(qm.Where("id_shows=?", showId)).All(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = showUrl.DeleteAll(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, "")

}

/*
	1. times_aired++
	2. last_aired = time.Now()
*/
func (server *Server) UpdateGoLive(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	show, err := models.Shows(qm.Where("id=?", id)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	/*update show*/
	show.TimesAired = null.Int{Int: show.TimesAired.Int + 1, Valid: true}
	show.LastAired = null.Time{Time: time.Now(), Valid: true}

	_, err = show.Update(context.Background(), server.DB, boil.Infer())
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, show)

}

func (server *Server) SetActiveShow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	show, err := models.Shows(qm.Where("id=?", id)).One(context.Background(), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	/*activate or deactivate show*/
	show.Active = !strings.Contains(r.URL.Path, "deactivate")

	_, err = show.Update(context.Background(), server.DB, boil.Infer())
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, show)

}

/* helper function */
func storeLogoImage(file multipart.File, handler *multipart.FileHeader) (string, error) {

	/* check file type is image */
	header := make([]byte, 512)
	if _, err := file.Read(header); err != nil {
		return "", err
	}
	if !strings.Contains(http.DetectContentType(header), "image") {
		return "", errors.New("input file is not an image")
	}

	defer file.Close()

	// TODO remove hard-coded path
	logoFileName := "/etc/logos/" + handler.Filename
	f, err := os.OpenFile(logoFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}

	return logoFileName, nil

}
