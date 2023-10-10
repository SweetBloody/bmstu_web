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

type teamHandler struct {
	teamUsecase models.TeamUsecaseI
}

func NewTeamHandler(m *mux.Router, teamUsecase models.TeamUsecaseI) {
	handler := &teamHandler{
		teamUsecase: teamUsecase,
	}

	m.Handle("/api/teams", middleware.AuthMiddleware(http.HandlerFunc(handler.GetTeamsOfSeason), "admin", "user")).Queries("season", "{season}").Methods("GET")
	m.Handle("/api/teams", middleware.AuthMiddleware(http.HandlerFunc(handler.GetAll), "admin", "user")).Methods("GET")
	m.Handle("/api/teams/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.GetTeamById), "admin", "user")).Methods("GET")
	//m.Handle("/api/teams_of_season/{season}", middleware.AuthMiddleware(http.HandlerFunc(handler.GetTeamsOfSeason), "admin", "user")).Methods("GET")
	m.Handle("/api/teams", middleware.AuthMiddleware(http.HandlerFunc(handler.Create), "admin")).Methods("POST")
	m.Handle("/api/teams/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Update), "admin")).Methods("PUT")
	m.Handle("/api/teams/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Delete), "admin")).Methods("DELETE")
}

// @Summary Get all teams
// @Tags teams
// @Description Get all teams
// @ID get-all-teams
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Team
// @Failure 500
// @Router /api/teams [get]
func (handler *teamHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	teams, err := handler.teamUsecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(teams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get team by id
// @Tags teams
// @Description Get team by id
// @ID get-team-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Team
// @Failure 500
// @Router /api/team/{id} [get]
func (handler *teamHandler) GetTeamById(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	team, err := handler.teamUsecase.GetTeamById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get all teams
// @Tags teams
// @Description Get all teams
// @ID get-all-teams
// @Accept  json
// @Produce  json
// @Param season query string false "season"
// @Success 200 {object} models.Team
// @Failure 500
// @Router /api/teams [get]
func (handler *teamHandler) GetTeamsOfSeason(w http.ResponseWriter, r *http.Request) {
	var season int
	var err error
	res := r.URL.Query().Get("season")
	season, err = strconv.Atoi(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	teams, err := handler.teamUsecase.GetTeamsOfSeason(season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(teams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Create team
// @Tags teams
// @Description Create team
// @ID create-team
// @Accept  json
// @Produce  json
// @Param input body models.Team true "team info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/teams [post]
func (handler *teamHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	team := new(models.Team)
	err := decoder.Decode(team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := handler.teamUsecase.Create(team)
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

// @Summary Update team
// @Tags teams
// @Description update teams
// @ID update-teams
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param input body models.Team true "team info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/teams/{id} [put]
func (handler *teamHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(r.Body)
	team := new(models.Team)
	err = decoder.Decode(team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.teamUsecase.Update(id, team)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete team
// @Tags teams
// @Description delete team
// @ID delete-team
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200
// @Failure 500
// @Router /api/teams/{id} [delete]
func (handler *teamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = handler.teamUsecase.Delete(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
