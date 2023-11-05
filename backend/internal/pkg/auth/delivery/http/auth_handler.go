package http

import (
	"encoding/json"
	"github.com/SweetBloody/bmstu_web/backend/internal/pkg/auth"
	"github.com/SweetBloody/bmstu_web/backend/internal/pkg/auth/token"
	"github.com/SweetBloody/bmstu_web/backend/internal/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
)

type authHandler struct {
	userUsecase models.UserUsecaseI
}

func NewAuthHandler(m *mux.Router, userUsecase models.UserUsecaseI) {
	handler := &authHandler{
		userUsecase: userUsecase,
	}
	m.HandleFunc("/api/login", handler.LogIn).Methods("POST")
	m.HandleFunc("/api/register", handler.Register).Methods("POST")
	m.HandleFunc("/api/logout", handler.LogOut).Methods("DELETE")
}

// @Summary Log in
// @Tags auth
// @Description User log in
// @ID auth-log-in
// @Accept  json
// @Produce  json
// @Param input body auth.LogInData true "account info"
// @Success 200 {object} string "Logged in successfully!"
// @Failure 400,401
// @Failure 500
// @Router /auth/login [post]
func (handler *authHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	loginData := new(auth.LogInData)
	err := decoder.Decode(loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ok, err := handler.userUsecase.Authenticate(loginData.Login, loginData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}
	user, err := handler.userUsecase.GetUserByLogin(loginData.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tokenString, err := token.GenerateToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:   "jwt-token",
		Value:  tokenString,
		MaxAge: 60 * 60 * 24,
		Path:   "/",
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("Logged in successfully!"))
}

// @Summary Register
// @Tags auth
// @Description New user register
// @ID auth-register
// @Accept  json
// @Produce  json
// @Param input body models.User true "account info"
// @Success 200
// @Failure 500
// @Router /auth/register [post]
func (handler *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := new(models.User)
	err := decoder.Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Role = "user"
	id, err := handler.userUsecase.Create(user)
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

// @Summary Log out
// @Tags auth
// @Description User log out
// @ID auth-log-out
// @Accept  json
// @Produce  json
// @Success 200 {object} string "Logged out successfully!"
// @Router /auth/logout [delete]
func (handler *authHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "jwt-token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("Logged out successfully!"))
}
