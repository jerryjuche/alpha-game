package user

import (
	"encoding/json"
	"net/http"

	"github.com/jerryjuche/alpha-game/internal/auth"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(auth.UserIDKey).(string)

	profile, err := u.service.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)

}
