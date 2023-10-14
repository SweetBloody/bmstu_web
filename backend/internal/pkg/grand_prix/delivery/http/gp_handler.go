package http

import (
	"encoding/json"
	"fmt"
	"github.com/SweetBloody/bmstu_web/backend/internal/app/middleware"
	"github.com/SweetBloody/bmstu_web/backend/internal/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type grandPrixHandler struct {
	grandPrixUsecase  models.GrandPrixUsecaseI
	raceResultUsecase models.RaceResultUsecaseI
	qualResultUsecase models.QualResultUsecaseI
}

func NewDriverHandler(m *mux.Router,
	grandPrixUsecase models.GrandPrixUsecaseI,
	raceResultUsecase models.RaceResultUsecaseI,
	qualResultUsecase models.QualResultUsecaseI) {
	handler := &grandPrixHandler{
		grandPrixUsecase:  grandPrixUsecase,
		raceResultUsecase: raceResultUsecase,
		qualResultUsecase: qualResultUsecase,
	}

	m.HandleFunc("/api/grandprix", handler.GetAllBySeason).Queries("season", "{season}").Methods("GET")
	m.HandleFunc("/api/grandprix", handler.GetAll).Methods("GET")
	//m.HandleFunc("/api/grandprix/id/{id}", handler.GetGPById).Methods("GET")
	m.HandleFunc("/api/grandprix/{id}", handler.GetGPById).Methods("GET")
	//m.HandleFunc("/api/grandprix/place/{place}", handler.GetAllByPlace).Methods("GET")
	m.Handle("/api/grandprix", middleware.AuthMiddleware(http.HandlerFunc(handler.Create), "admin")).Methods("POST")
	m.Handle("/api/grandprix/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Update), "admin")).Methods("PUT")
	m.Handle("/api/grandprix/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Delete), "admin")).Methods("DELETE")
	m.Handle("/api/grandprix/{id}/name", middleware.AuthMiddleware(http.HandlerFunc(handler.UpdateGPName), "admin")).Methods("PATCH")
	m.Handle("/api/grandprix/{id}/race_results", middleware.AuthMiddleware(http.HandlerFunc(handler.GetRaceResultsOfGP), "admin", "user")).Methods("GET")
	m.Handle("/api/grandprix/{id}/qual_results", middleware.AuthMiddleware(http.HandlerFunc(handler.GetQualResultsOfGP), "admin", "user")).Methods("GET")
}

// @Summary Get all gp
// @Tags gp
// @Description Get all gp
// @ID get-all-gp
// @Accept  json
// @Produce  json
// @Success 200 {object} models.GrandPrix
// @Failure 500
// @Router /api/grandprix [get]
func (handler *grandPrixHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	gp, err := handler.grandPrixUsecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(gp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get gp by id
// @Tags gp
// @Description Get gp by id
// @ID get-gp-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.GrandPrix
// @Failure 400
// @Failure 500
// @Router /api/grandprix/{id} [get]
func (hander *grandPrixHandler) GetGPById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	gp, err := hander.grandPrixUsecase.GetGPById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(gp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get all gp
// @Tags gp
// @Description Get all gp
// @ID get-all-gp
// @Accept  json
// @Produce  json
// @Param season query string false "season"
// @Success 200 {object} models.GrandPrix
// @Failure 500
// @Router /api/grandprix [get]
func (hander *grandPrixHandler) GetAllBySeason(w http.ResponseWriter, r *http.Request) {
	res := r.URL.Query().Get("season")
	season, err := strconv.Atoi(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	gp, err := hander.grandPrixUsecase.GetAllBySeason(season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(gp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//func (hander *grandPrixHandler) GetAllByPlace(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	place := vars["place"]
//	encoder := json.NewEncoder(w)
//	gp, err := hander.grandPrixUsecase.GetAllByPlace(place)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	err = encoder.Encode(gp)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

// @Summary Create gp
// @Tags gp
// @Description Create gp
// @ID create-gp
// @Accept  json
// @Produce  json
// @Param input body models.GrandPrix true "GP info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/grandprix [post]
func (handler *grandPrixHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	gp := new(models.GrandPrix)
	err := decoder.Decode(gp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := handler.grandPrixUsecase.Create(gp)
	if err != nil {
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

// @Summary Update gp
// @Tags gp
// @Description Update gp
// @ID update-gp
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param input body models.GrandPrix true "GP info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/grandprix/{id} [put]
func (handler *grandPrixHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(r.Body)
	gp := new(models.GrandPrix)
	err = decoder.Decode(gp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.grandPrixUsecase.Update(id, gp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete gp
// @Tags gp
// @Description delete gp
// @ID delete-gp
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/grandprix/{id} [delete]
func (handler *grandPrixHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.grandPrixUsecase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Update gp name
// @Tags gp
// @Description Update gp name
// @ID update-gp-name
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param input body string true "gp_name"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/grandprix/{id}/name [patch]
func (handler *grandPrixHandler) UpdateGPName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(r.Body)
	req := struct {
		Name string `json:"gp_name"`
	}{}
	err = decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.grandPrixUsecase.UpdateGPName(id, req.Name)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get raceresults of gp
// @Tags race_results
// @Description Get raceresults of gp
// @ID get-race-of-gp
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.RaceResultView
// @Failure 400
// @Failure 500
// @Router /api/grandprix/{id}/race_results [get]
func (handler *grandPrixHandler) GetRaceResultsOfGP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	raceResults, err := handler.raceResultUsecase.GetRaceResultsOfGP(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(raceResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get qualresults of gp
// @Tags qual_results
// @Description Get qualresults of gp
// @ID get-qual-of-gp
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.QualResultView
// @Failure 400
// @Failure 500
// @Router /api/grandprix/{id}/qual_results [get]
func (handler *grandPrixHandler) GetQualResultsOfGP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	qualResults, err := handler.qualResultUsecase.GetQualResultsOfGP(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(qualResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
