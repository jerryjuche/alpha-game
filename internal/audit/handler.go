package audit

import (
	"encoding/json"
	"net/http"
)

type AuditHandler struct {
	service *AuditService
}

func NewAuditHandler(service *AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (s *AuditHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	pendingSubmissions, err := s.service.GetPendingSubmissions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pendingSubmissions)
}

func (s *AuditHandler) Approve(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SubmissionID string `json:"submission_id"`
		Points       int    `json:"points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	auditorID := r.Context().Value("user_id").(string)

	if err := s.service.ApproveSubmission(r.Context(), input.SubmissionID, auditorID, input.Points); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "approved",
	})
}

func (s *AuditHandler) Reject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SubmissionID string `json:"submission_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	auditorID := r.Context().Value("user_id").(string)

	if err := s.service.RejectSubmission(r.Context(), input.SubmissionID, auditorID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "rejected",
	})
}