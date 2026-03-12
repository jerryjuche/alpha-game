package scoring

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ScoringService struct {
	DBConn *sqlx.DB
}

func NewScoringService(db *sqlx.DB) *ScoringService {
	return &ScoringService{DBConn: db}
}

func (g *ScoringService) AwardPoints(ctx context.Context, playerID string, gameID string, points int) error {
	_, err := g.DBConn.ExecContext(ctx, "UPDATE game_players SET score = score + $1 WHERE user_id = $2 AND game_id = $3 ", points, playerID, gameID)
	if err != nil {
		return fmt.Errorf("error updating scores: %w", err)
	}
	return nil
}
