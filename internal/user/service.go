package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserService struct {
	DBConn *sqlx.DB
}

type UserProfile struct {
	Level        int
	GamesPlayed  int
	GamesWon     int
	WinRate      float64
	LongestWord  string
	ShortestWord string
	TotalCorrect int
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{DBConn: db}

}

func (a *UserService) GetProfile(ctx context.Context, userID string) (*UserProfile, error) {
	var profile UserProfile

	err := a.DBConn.QueryRowContext(ctx, "SELECT level FROM users WHERE id = $1", userID).Scan(&profile.Level)
	if err != nil {
		return nil, fmt.Errorf("Error getting user stats, %w", err)
	}

	err = a.DBConn.QueryRowContext(ctx, "SELECT COUNT (*) FROM game_players WHERE id = $1", userID).Scan(&profile.GamesPlayed)
	if err != nil {
		return nil, fmt.Errorf("Error getting user stats, %w", err)
	}

	err = a.DBConn.QueryRowContext(ctx, "SELECT COUNT (*) FROM submissions WHERE submitted_by = $1 AND status = 'approved'", userID).Scan(&profile.TotalCorrect)
	if err != nil {
		return nil, fmt.Errorf("Error getting user stats, %w", err)
	}

	err = a.DBConn.QueryRowContext(ctx, "SELECT MAX(word_submitted), MIN(word_submitted) FROM submissions WHERE submitted_by = $1 AND status = 'approved'", userID).Scan(&profile.LongestWord, &profile.ShortestWord)
	if err != nil {
		return nil, fmt.Errorf("Error getting user stats")
	}

	if profile.GamesPlayed > 0 {
		profile.WinRate = float64(profile.GamesWon) / float64(profile.GamesPlayed) * 100
	}
	return &profile, nil

}
