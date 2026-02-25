package game

import (
	"encoding/json"
	"net/http"

	"github.com/jerryjuche/alpha-game/internal/auth"
)

type GameHandler struct {
	service *GameEngine
}

func NewGameHandler(service *GameEngine) *GameHandler {
	return &GameHandler{service: service}
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	hostID := r.Context().Value(auth.UserIDKey).(string)

	gameID, inviteCode, err := h.service.CreateGame(r.Context(), hostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"game_id":     gameID,
		"invite_code": inviteCode,
	})
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	playerID := r.Context().Value(auth.UserIDKey).(string)

	var input struct {
		InviteCode string `json:"invite_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameID, err := h.service.JoinGame(r.Context(), playerID, input.InviteCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"game_id": gameID,
	})
}

func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	hostID := r.Context().Value(auth.UserIDKey).(string)

	var input struct {
		GameID string `json:"game_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameID, err := h.service.StartGame(r.Context(), input.GameID, hostID, "active")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"game_id": gameID,
	})
}

func (h *GameHandler) Submission(w http.ResponseWriter, r *http.Request) {
	playerID := r.Context().Value(auth.UserIDKey).(string)

	var input struct {
		RoundID  string `json:"round_id"`
		Word     string `json:"word"`
		Category string `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, err := h.service.SubmitAnswer(r.Context(), "", playerID, input.RoundID, input.Category, input.Word)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if valid {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "approved",
			"word":     input.Word,
			"category": input.Category,
			"word_id":  input.RoundID,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "rejected",
			"word":     input.Word,
			"category": input.Category,
			"word_id":  input.RoundID,
		})
	}

}
