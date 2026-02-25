package audit

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuditService struct {
	DBConn *sqlx.DB
}

type PendingSubmission struct {
	ID          string
	SubmittedBy string
	Word        string
	Category    string
}

func NewAuditService(db *sqlx.DB) *AuditService {
	return &AuditService{DBConn: db}
}

func (a *AuditService) GetPendingSubmissions(ctx context.Context) ([]PendingSubmission, error) {

	rows, err := a.DBConn.QueryContext(ctx, "SELECT id, submitted_by, category, word_submitted FROM submissions WHERE status = 'pending'")

	if err != nil {
		return nil, fmt.Errorf("error cfetching submissions, %w", err)
	}

	var submissions []PendingSubmission
	for rows.Next() {
		var s PendingSubmission
		if err := rows.Scan(&s.ID, &s.SubmittedBy, &s.Category, &s.Word); err != nil {
			return nil, fmt.Errorf("error scanning rows, %w", err)
		}
		submissions = append(submissions, s)
	}
	return submissions, nil

}

func (a *AuditService) ApproveSubmission(ctx context.Context) error {
	var submission_id string
	submissionID, err := a.DBConn.QueryContext(ctx, "SELECT id FROM submissions WHERE status = 'pending'").Scan(&submission_id)
	if err != nil {
		return fmt.Errorf("error fetching rows, %w", err)
	}


}
