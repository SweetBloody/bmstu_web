package http

import (
	"encoding/json"
	"fmt"
	"github.com/SweetBloody/bmstu_web/backend/internal/app/middleware"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/SweetBloody/bmstu_web/backend/internal/pkg/models"
)

type qualResultHandler struct {
	qualResultUsecase models.QualResultUsecaseI
}

func NewQualResultHandler(m *mux.Router, qualResultUsecase models.QualResultUsecaseI) {
	handler := &qualResultHandler{
		qualResultUsecase: qualResultUsecase,
	}

	m.Handle("/api/grandprix/{gp_id}/qual_results", middleware.AuthMiddleware(http.HandlerFunc(handler.Create), "admin")).Methods("POST")
	//m.Handle("/api/qual_results", middleware.AuthMiddleware(http.HandlerFunc(handler.GetAll), "admin", "user")).Methods("GET")
	//m.Handle("/api/grandprix/{gp_id}/qual_results/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.GetQualResultById), "admin", "user")).Methods("GET")
	m.Handle("/api/grandprix/{gp_id}/qual_results/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Update), "admin")).Methods("PUT")
	m.Handle("/api/grandprix/{gp_id}/qual_results/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Delete), "admin")).Methods("DELETE")
}

//func (handler *qualResultHandler) GetAll(w http.ResponseWriter, r *http.Request) {
//	encoder := json.NewEncoder(w)
//	results, err := handler.qualResultUsecase.GetAll()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	err = encoder.Encode(results)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

func (handler *qualResultHandler) GetQualResultById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	result, err := handler.qualResultUsecase.GetQualResultById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = encoder.Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Create qual_result
// @Tags qual_results
// @Description Create qual_result
// @ID create-qual_result
// @Accept  json
// @Produce  json
// @Param gp_id path string true "gp_id"
// @Param input body models.QualResult true "qual result info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/grandprix/{gp_id}/qual_results [post]
func (handler *qualResultHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	result := new(models.QualResult)
	err := decoder.Decode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := handler.qualResultUsecase.Create(result)
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

// @Summary Update qual_result
// @Tags qual_results
// @Description Update qual_result
// @ID update-qual_result
// @Accept  json
// @Produce  json
// @Param gp_id path string true "gp_id"
// @Param id path string true "id"
// @Param input body models.QualResult true "qual result info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/grandprix/{gp_id}/qual_results/{id} [put]
func (handler *qualResultHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(r.Body)
	result := new(models.QualResult)
	err = decoder.Decode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.qualResultUsecase.Update(id, result)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete qual_result
// @Tags qual_results
// @Description delete qual_result
// @ID delete-qual_result
// @Accept  json
// @Produce  json
// @Param gp_id path string true "gp_id"
// @Param id path string true "id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/grandprix/{gp_id}/qual_results/{id} [delete]
func (handler *qualResultHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = handler.qualResultUsecase.Delete(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
