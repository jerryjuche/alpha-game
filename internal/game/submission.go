package game

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var ErrDuplicateSubmission = errors.New("submission already exists")

func (g *GameEngine) SubmitAnswer(ctx context.Context, gameID string, playerID string, roundID string, category string, word string) (bool, error) {

	_, err := g.DBConn.ExecContext(ctx, "INSERT INTO submissions (round_id, submitted_by, category, word_submitted, status) VALUES ($1, $2, $3, $4, $5)", roundID, playerID, category, word, "pending")
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return false, ErrDuplicateSubmission
		}
		return false, fmt.Errorf("cant insert into db, %w", err)
	}

	valid, err := g.WordService.LookupWord(ctx, word, category)
	if err != nil {
		return false, fmt.Errorf("error submiting, %w", err)
	}

	status := "pending"

	if valid {
		status = "approved"
		err = g.Scoring.AwardPoints(ctx, playerID, gameID, 10)
		if err != nil {
			return false, fmt.Errorf("error awarding points: %w", err)
		}
	} 

	_, err = g.DBConn.ExecContext(ctx, "UPDATE submissions SET status = $1 WHERE round_id = $2 AND submitted_by = $3", status, roundID, playerID)
	if err != nil {
		return false, fmt.Errorf("Error updating the db, %w", err)
	}

	return valid, nil

}
