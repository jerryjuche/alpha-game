package word

import (
	"encoding/json"
	"net/http"

	"github.com/jerryjuche/alpha-game/internal/auth"
)

type WordHandler struct {
	service *WordService
}

func NewWordHandler(service *WordService) *WordHandler {
	return &WordHandler{service: service}
}

func (h *WordHandler) AddWord(w http.ResponseWriter, r *http.Request) {

	var input AddWordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error adding word to db", http.StatusBadRequest)
		return
	}

	if err := h.service.AddWord(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&input)

}

func (h *WordHandler) DeleteWord(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WordID string `json:"word_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteWord(r.Context(), input.WordID); err != nil {
		http.Error(w, "error deleting word", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&input)

}

func (h *WordHandler) ApproveWord(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WordID string `json:"word_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ApproveWord(r.Context(), input.WordID); err != nil {
		http.Error(w, "error approving word", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&input)

}

func (h *WordHandler) ImportFromExcel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	category := r.FormValue("category")

	userID := r.Context().Value(auth.UserIDKey).(string)

	if _, err := h.service.ImportFromExcel(r.Context(), file, category, userID); err != nil {
		http.Error(w, "error exporting file", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "imported"})

}
