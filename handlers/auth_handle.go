package handlers

import (
	"cloneRaddit/user"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	Repo user.Repo
}

func NewUserHandler(repo user.Repo) *UserHandler {
	return &UserHandler{Repo: repo}
}

type Req struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

var jwtSecret = []byte("SecretKey")

type jwtResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req = Req{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Невалидный Json"+err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.Repo.Register(req.Name, req.Password)
	if err != nil {
		http.Error(w, "пользователь уже зарегистрирован"+err.Error(), http.StatusConflict)
		return
	}

	token, err := GenerateJWTToken(u)
	if err != nil {
		http.Error(w, "Не удалось получить токен"+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(jwtResponse{Token: token})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req = Req{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Невалидный Json"+err.Error(), http.StatusBadRequest)
		return
	}
	u, err := h.Repo.Login(req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	token, err := GenerateJWTToken(u)
	if err != nil {
		http.Error(w, "Не удалось получить токен"+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(jwtResponse{Token: token})
}

func GenerateJWTToken(u *user.User) (string, error) {
	claims := jwt.MapClaims{
		"user": map[string]string{
			"id":       u.ID,
			"username": u.Name,
		},
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
