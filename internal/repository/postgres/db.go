package internal

import (
	"fmt"

	"github.com/jerryjuche/alpha-game/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(cfg *config.Config) (*sqlx.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		
	)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error opening database, %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error conecting to database, %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil

}
