package http

import (
	"encoding/json"
	"fmt"
	"github.com/SweetBloody/bmstu_web/backend/internal/app/middleware"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"

	"github.com/SweetBloody/bmstu_web/backend/internal/pkg/models"
)

type driverHandler struct {
	driverUsecase     models.DriverUsecaseI
	raceResultUsecase models.RaceResultUsecaseI
}

func NewDriverHandler(m *mux.Router, driverUsecase models.DriverUsecaseI, raceResultUsecase models.RaceResultUsecaseI) {
	handler := &driverHandler{
		driverUsecase:     driverUsecase,
		raceResultUsecase: raceResultUsecase,
	}

	m.HandleFunc("/api/drivers", handler.GetDriversOfSeason).Queries("season", "{season}").Methods("GET")
	m.Handle("/api/drivers", middleware.AuthMiddleware(http.HandlerFunc(handler.GetAll), "admin", "user")).Methods("GET")
	m.Handle("/api/drivers", middleware.AuthMiddleware(http.HandlerFunc(handler.Create), "admin")).Methods("POST")
	m.Handle("/api/drivers/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.GetDriverById), "admin", "user")).Methods("GET")
	//m.Handle("/api/drivers", middleware.AuthMiddleware(http.HandlerFunc(handler.GetDriversOfSeason), "admin", "user")).Queries("season", "{season}").Methods("GET")
	m.Handle("/api/drivers_standing", middleware.AuthMiddleware(http.HandlerFunc(handler.GetDriversStanding), "admin", "user")).Methods("GET")
	m.Handle("/api/drivers/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Update), "admin")).Methods("PUT")
	m.Handle("/api/drivers/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Delete), "admin")).Methods("DELETE")
	m.Handle("/api/drivers_teams", middleware.AuthMiddleware(http.HandlerFunc(handler.LinkDriverTeam), "admin")).Methods("POST")
	m.HandleFunc("/api/drivers", handler.GetRaceWinnerOfGP).Queries("winner_gp_id", "{winner_gp_id}").Methods("GET")
}

// @Summary Get all drivers
// @Tags drivers
// @Description Get all drivers
// @ID get-all-drivers
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Driver
// @Failure 500
// @Router /api/drivers [get]
func (handler *driverHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	drivers, err := handler.driverUsecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(drivers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get driver by id
// @Tags drivers
// @Description Get driver by id
// @ID get-driver-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Driver
// @Failure 400
// @Failure 500
// @Router /api/drivers/{id} [get]
func (handler *driverHandler) GetDriverById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	driver, err := handler.driverUsecase.GetDriverById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get all drivers
// @Tags drivers
// @Description Get all drivers
// @ID get-all-drivers
// @Accept  json
// @Produce  json
// @Param season query string false "season"
// @Success 200 {object} models.Driver
// @Failure 500
// @Router /api/drivers [get]
func (handler *driverHandler) GetDriversOfSeason(w http.ResponseWriter, r *http.Request) {
	var season int
	var err error
	_, err = r.Cookie("jwt-token")
	if err != nil {
		season = time.Now().Year() - 1
	} else {
		res := r.URL.Query().Get("season")
		season, err = strconv.Atoi(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	encoder := json.NewEncoder(w)
	drivers, err := handler.driverUsecase.GetDriversOfSeason(season)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(drivers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get drivers standings
// @Tags drivers
// @Description Get drivers standings
// @ID get-drivers-standings
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Standings
// @Failure 500
// @Router /api/drivers_standing [get]
func (handler *driverHandler) GetDriversStanding(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	standing, err := handler.driverUsecase.GetDriversStanding()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(standing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Create driver
// @Tags drivers
// @Description Create driver
// @ID create-driver
// @Accept  json
// @Produce  json
// @Param input body models.Driver true "driver info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/drivers [post]
func (handler *driverHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	driver := new(models.Driver)
	err := decoder.Decode(driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := handler.driverUsecase.Create(driver)
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

// @Summary Update driver
// @Tags drivers
// @Description Update driver
// @ID update-driver
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param input body models.Driver true "driver info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/drivers/{id} [put]
func (handler *driverHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(r.Body)
	driver := new(models.Driver)
	err = decoder.Decode(driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.driverUsecase.Update(id, driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete driver
// @Tags drivers
// @Description delete driver
// @ID delete-driver
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/drivers/{id} [delete]
func (handler *driverHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = handler.driverUsecase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Link driver team
// @Tags drivers
// @Description Link driver team
// @ID link-driver-team
// @Accept  json
// @Produce  json
// @Param input body models.DriversTeams true "driver-team connection info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/drivers_teams [post]
func (handler *driverHandler) LinkDriverTeam(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	link := new(models.DriversTeams)
	err := decoder.Decode(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.driverUsecase.LinkDriverTeam(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get race winner of gp
// @Tags drivers
// @Description Get race winner of gp
// @ID get-race-winner-of-gp
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param winner_gp_id query string false "winner_gp_id"
// @Success 200 {object} models.RaceResultView
// @Failure 400
// @Failure 500
// @Router /api/drivers [get]
func (handler *driverHandler) GetRaceWinnerOfGP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	qualResults, err := handler.raceResultUsecase.GetRaceWinnerOfGP(id)
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
