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

type userHandler struct {
	userUsecase models.UserUsecaseI
}

func NewUserHandler(m *mux.Router, userUsecase models.UserUsecaseI) {
	handler := &userHandler{
		userUsecase: userUsecase,
	}

	m.Handle("/api/users", middleware.AuthMiddleware(http.HandlerFunc(handler.Create), "admin")).Methods("POST")
	m.Handle("/api/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Update), "admin")).Methods("PUT")
	m.Handle("/api/users/{id}", middleware.AuthMiddleware(http.HandlerFunc(handler.Delete), "admin")).Methods("DELETE")
}

// @Summary Create user
// @Tags users
// @Description create user
// @ID create-user
// @Accept  json
// @Produce  json
// @Param input body models.User true "user info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/users [post]
func (handler *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := new(models.User)
	err := decoder.Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := handler.userUsecase.Create(user)
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

// @Summary Update user
// @Tags users
// @Description update user
// @ID update-user
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param input body models.User true "user info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/users/{id} [put]
func (handler *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(r.Body)
	user := new(models.User)
	err = decoder.Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.userUsecase.Update(id, user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete user
// @Tags users
// @Description delete user
// @ID delete-user
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200
// @Failure 500
// @Router /api/users/{id} [delete]
func (handler *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = handler.userUsecase.Delete(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
