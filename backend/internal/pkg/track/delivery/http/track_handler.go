package http

import (
	"encoding/json"
	"fmt"
	middleware "github.com/SweetBloody/bmstu_web/backend/internal/app/middleware"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/SweetBloody/bmstu_web/backend/internal/pkg/models"
)

type trackHandler struct {
	trackUsecase models.TrackUsecaseI
}

func NewTrackHandler(m *mux.Router, trackUsecase models.TrackUsecaseI) {
	handler := &trackHandler{
		trackUsecase: trackUsecase,
	}

	m.Handle("/api/tracks", middleware.AuthMiddleware(http.HandlerFunc(handler.Create), "admin")).Methods("POST")
	m.Handle("/api/tracks", middleware.AuthMiddleware(http.HandlerFunc(handler.GetAll), "admin", "user")).Methods("GET")
	m.Handle("/api/tracks/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.GetDriverById), "admin", "user")).Methods("GET")
	m.Handle("/api/tracks/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Update), "admin")).Methods("PUT")
	m.Handle("/api/tracks/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Delete), "admin")).Methods("DELETE")
}

// @Summary Get all tracks
// @Tags tracks
// @Description Get all tracks
// @ID get-all-tracks
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Track
// @Failure 500
// @Router /api/tracks [get]
func (handler *trackHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	tracks, err := handler.trackUsecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(tracks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get track by id
// @Tags tracks
// @Description Get track by id
// @ID get-track-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Track
// @Failure 500
// @Router /api/tracks/{id} [get]
func (handler *trackHandler) GetDriverById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	track, err := handler.trackUsecase.GetTeamById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Create track
// @Tags tracks
// @Description create track
// @ID create-track
// @Accept  json
// @Produce  json
// @Param input body models.Track true "track info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/tracks [post]
func (handler *trackHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	track := new(models.Track)
	err := decoder.Decode(track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := handler.trackUsecase.Create(track)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Update track
// @Tags tracks
// @Description update track
// @ID update-track
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param input body models.Track true "track info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/tracks/{id} [put]
func (handler *trackHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(r.Body)
	track := new(models.Track)
	err = decoder.Decode(track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.trackUsecase.Update(id, track)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete track
// @Tags tracks
// @Description delete track
// @ID delete-track
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200
// @Failure 500
// @Router /api/tracks/{id} [delete]
func (handler *trackHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = handler.trackUsecase.Delete(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
