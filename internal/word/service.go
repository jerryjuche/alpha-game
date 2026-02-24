package word

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type WordService struct {
	DBConn *sqlx.DB
}

type AddWordInput struct {
	Word     string
	Category string
	AddedBy  string
}

func NewWordService(db *sqlx.DB) *WordService {
	return &WordService{
		DBConn: db,
	}
}

func (w *WordService) AddWord(ctx context.Context, input AddWordInput) error {

	firstLetter := strings.ToUpper(string(input.Word[0]))

	_, err := w.DBConn.ExecContext(ctx, "INSERT INTO word_database (word, category, starts_with, added_by) VALUES ($1, $2, $3, $4)", input.Word, input.Category, firstLetter, input.AddedBy)
	if err != nil {
		return fmt.Errorf("Error writing to DB, %w", err)
	}
	return nil

}

func (w *WordService) ApproveWord(ctx context.Context, wordID string) error {

	_, err := w.DBConn.ExecContext(ctx, "UPDATE word_database SET is_approved = true WHERE id = $1", wordID)
	if err != nil {
		return fmt.Errorf("Cannot approve word, %w", err)
	}
	return nil

}

func (w *WordService) DeleteWord(ctx context.Context, wordID string) error {
	_, err := w.DBConn.ExecContext(ctx, "DELETE FROM word_database WHERE id = $1", wordID)
	if err != nil {
		return fmt.Errorf("Cannot Delete %v from DB, %w", wordID, err)
	}
	return nil
}

func (w *WordService) LookupWord(ctx context.Context, word string, category string, letter string) (bool, error) {
	var verify string
	letter = strings.ToUpper(string(word[0]))

	err := w.DBConn.QueryRowContext(ctx, "SELECT id FROM word_database WHERE word = $1 AND category = $2 AND starts_with = $3 AND is_approved = true", word, category, letter).Scan(&verify)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("lookup error, %w", err)
	}

	return true, nil
}
