package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DBConn    *sqlx.DB
	JWTSecret string
}

type RegisterInput struct {
	Username string
	Password string
	Email    string
}

type LoginInput struct {
	Email    string
	Password string
}


func NewAuthService(db *sqlx.DB, jwtSecret string) *AuthService {
	return &AuthService{
		DBConn:    db,
		JWTSecret: jwtSecret,
	}
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (string, error) {
	// Check if email already exists
	var userID string
	var hashedPassword string

	err := s.DBConn.QueryRowxContext(ctx, "SELECT id, password FROM users WHERE email = $1", input.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", fmt.Errorf("User not found!!: %w", err)
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		return "", fmt.Errorf("Incorrect password: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(userID)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return token, nil

}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (string, error) {
	// Check if email already exists
	var existingID string
	err := s.DBConn.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", input.Email).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error checking email: %w", err)
	}
	if existingID != "" {
		return "", fmt.Errorf("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	// Save user to database
	var userID string
	err = s.DBConn.QueryRowContext(ctx,
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		input.Username, input.Email, string(hashedPassword),
	).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(userID)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return token, nil
}

func (s *AuthService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}
