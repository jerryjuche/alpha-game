package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jerryjuche/alpha-game/config"
	"github.com/jerryjuche/alpha-game/internal/audit"
	"github.com/jerryjuche/alpha-game/internal/auth"
	"github.com/jerryjuche/alpha-game/internal/game"
	pg "github.com/jerryjuche/alpha-game/internal/repository/postgres"
	ws "github.com/jerryjuche/alpha-game/internal/websocket"
	"github.com/jerryjuche/alpha-game/internal/word"
)

func main() {
	cfg := config.Load()
	fmt.Printf("Server starting on port %d in %s mode\n", cfg.AppPort, cfg.Env)

	db, err := pg.NewDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Database connected successfully!")

	// Auth

	hub := ws.NewHub()
	go hub.Run()

	auditService := audit.NewAuditService(db)
	auditHandler := audit.NewAuditHandler(auditService)
	wordService := word.NewWordService(db)
	authService := auth.NewAuthService(db, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)
	gameEngine := game.NewGameEngine(db, hub, wordService)
	gameHandler := game.NewGameHandler(gameEngine)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/audit/pending", auditHandler.GetPending)
	r.Post("/audit/approve", auditHandler.Approve)
	r.Post("/audit/reject", auditHandler.Reject)
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/game/submit", gameHandler.Submission)

	r.Group(func(r chi.Router) {
		r.Use(authService.Authenticate)
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(auth.UserIDKey).(string)
			w.Write([]byte("Hello user: " + userID))
		})
		r.Post("/game/create", gameHandler.CreateGame)
		r.Post("/game/join", gameHandler.JoinGame)
		r.Post("/game/start", gameHandler.StartGame)
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.AppPort)
	log.Printf("Server running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))

}
