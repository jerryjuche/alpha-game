package audit

import (
	"context"
	"fmt"

	"github.com/jerryjuche/alpha-game/internal/scoring"
	"github.com/jmoiron/sqlx"
)

type AuditService struct {
	DBConn  *sqlx.DB
	Scoring *scoring.ScoringService
}

type PendingSubmission struct {
	ID          string
	SubmittedBy string
	Word        string
	Category    string
}

func NewAuditService(db *sqlx.DB, s *scoring.ScoringService) *AuditService {
	return &AuditService{
		DBConn:  db,
		Scoring: s,
	}
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

func (a *AuditService) ApproveSubmission(ctx context.Context, category string, input string, submissionID string, firstletter string, auditorID string, points int) error {
	_, err := a.DBConn.ExecContext(ctx, "UPDATE submissions SET status = 'approved', points_awarded = $1 WHERE id = $2", points, submissionID)
	if err != nil {
		return fmt.Errorf("error approving submissions, %w", err)
	}

	var playerID, gameID string

	err = a.DBConn.QueryRowContext(ctx, "SELECT s.submitted_by, r.game_id FROM submissions s JOIN rounds r ON r.id = s.round_id WHERE s.id = $1", submissionID).Scan(&playerID, &gameID)
	if err != nil {
		return fmt.Errorf("error fetching submission details: %w", err)
	}

	err = a.Scoring.AwardPoints(ctx, playerID, gameID, points)
	if err != nil {
		return fmt.Errorf("error awarding points: %w", err)
	}

	_, err = a.DBConn.ExecContext(ctx, "INSERT INTO audit_log (submission_id, reviewed_by, decision) VALUES ($1, $2, $3)", submissionID, auditorID, "approved")
	if err != nil {
		return fmt.Errorf("Error inserting to audit_log, %w", err)
	}

	_, err = a.DBConn.ExecContext(ctx, "INSERT INTO word_database (word, category, starts_with, is_approved) VALUES ($1, $2, $3, $4) ON CONFLICT (word, category) DO NOTHING", input, category, firstletter, true)
	if err != nil {
		return fmt.Errorf("error inserting into word_db", err)
	}

	return nil
}

func (a *AuditService) RejectSubmission(ctx context.Context, submissionID string, auditorID string) error {
	_, err := a.DBConn.ExecContext(ctx, "UPDATE submissions SET status = 'rejected', points_awarded = '0' WHERE id = $1", submissionID)
	if err != nil {
		return fmt.Errorf("Error rejecting submissions, %w", err)
	}

	_, err = a.DBConn.ExecContext(ctx, "INSERT INTO audit_log (submission_id, reviewed_by, decision) VALUES ($1, $2, $3)", submissionID, auditorID, "rejected")
	if err != nil {
		return fmt.Errorf("error inserting into audit log, %w", err)
	}

	return nil
}
