package game

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/jerryjuche/alpha-game/internal/scoring"
	ws "github.com/jerryjuche/alpha-game/internal/websocket"
	"github.com/jerryjuche/alpha-game/internal/word"
	"github.com/jmoiron/sqlx"
)

// Game holds all in-memory state for a single active game.
// mu protects all fields on this struct — any goroutine
// that reads OR writes these fields must hold the appropriate lock.
type Game struct {
	mu             sync.RWMutex
	Players        map[string]*Player
	GameId         string
	LetterUsage    map[string]int
	Start          time.Time
	End            time.Time
	Status         string
	CurrentLetter  string
	HostId         string
	PhaseEndsAt    time.Time
	RoundEndsAt    time.Time
	CurrentPhase   string
	CurrentRoundID string // NEW: tracked so reconnecting players can submit
}

type Player struct {
	Id           string
	Hints        int
	Score        int
	IsEliminated bool
}

// GameEngine manages all active games.
// mu protects the ActiveGames map itself — not the game objects inside it.
// Game objects have their own mu for field-level protection.
type GameEngine struct {
	mu          sync.RWMutex // guards ActiveGames map
	DBConn      *sqlx.DB
	Hub         *ws.Hub
	ActiveGames map[string]*Game
	WordService *word.WordService
	Scoring     *scoring.ScoringService
}

func NewGameEngine(db *sqlx.DB, h *ws.Hub, w *word.WordService, s *scoring.ScoringService) *GameEngine {
	return &GameEngine{
		DBConn:      db,
		Hub:         h,
		ActiveGames: make(map[string]*Game),
		WordService: w,
		Scoring:     s,
	}
}

func generateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	code := make([]byte, 6)
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[n.Int64()]
	}
	return string(code)
}

// loadGameFromDB fetches a game from the database and places it in ActiveGames.
// IMPORTANT: caller must already hold the engine write lock (mu.Lock).
// This is a private method — it is never called from outside this file.
func (g *GameEngine) loadGameFromDB(ctx context.Context, gameID string) error {
	var hostID string
	var status string
	err := g.DBConn.QueryRowContext(ctx,
		"SELECT created_by, status FROM games WHERE id = $1", gameID,
	).Scan(&hostID, &status)
	if err != nil {
		return fmt.Errorf("game not found in database: %w", err)
	}
	g.ActiveGames[gameID] = &Game{
		GameId:      gameID,
		HostId:      hostID,
		Status:      status,
		Players:     make(map[string]*Player),
		LetterUsage: make(map[string]int),
		Start:       time.Now(),
		End:         time.Now(),
	}
	return nil
}

// getGame retrieves a game from memory or loads it from the database.
//
// Why Lock() and not RLock() here?
// This is the check-then-act pattern:
//  1. Check if game exists in map  (read)
//  2. If not, load from DB         (write to map)
//
// Between step 1 and step 2, another goroutine could do the same check
// and also see nil — causing two goroutines to both write to the map.
// A full Lock() prevents this. There is no safe way to do this with RLock().
func (g *GameEngine) getGame(ctx context.Context, gameID string) (*Game, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	game := g.ActiveGames[gameID]
	if game == nil {
		if err := g.loadGameFromDB(ctx, gameID); err != nil {
			return nil, fmt.Errorf("game not found: %w", err)
		}
		game = g.ActiveGames[gameID]
	}
	return game, nil
}

// CreateGame inserts a new game into the database and registers it in memory.
// Always writes to the map → always Lock().
func (g *GameEngine) CreateGame(ctx context.Context, hostID string) (string, string, error) {
	var gameID string
	inviteCode := generateInviteCode()
	err := g.DBConn.QueryRowContext(ctx,
		"INSERT INTO games (invite_code, status, created_by) VALUES ($1, $2, $3) RETURNING id",
		inviteCode, "waiting", hostID,
	).Scan(&gameID)
	if err != nil {
		return "", "", fmt.Errorf("cannot create game: %w", err)
	}

	g.mu.Lock()
	g.ActiveGames[gameID] = &Game{
		GameId:        gameID,
		LetterUsage:   make(map[string]int),
		Players:       make(map[string]*Player),
		Start:         time.Now(),
		End:           time.Now(),
		Status:        "waiting",
		CurrentLetter: "",
		HostId:        hostID,
	}
	g.mu.Unlock()

	return gameID, inviteCode, nil
}

// JoinGame adds a player to an existing game.
//
// Two-level locking pattern:
//   - Engine lock (Lock) to safely retrieve the game from the map
//   - Game lock (Lock) to safely mutate game.Players
//
// Notice we release the engine lock before acquiring the game lock.
// Holding both simultaneously would risk deadlock if another goroutine
// acquires them in a different order.
func (g *GameEngine) JoinGame(ctx context.Context, playerID string, inviteCode string) (string, error) {
	var gameID string
	var status string
	err := g.DBConn.QueryRowContext(ctx,
		"SELECT id, status FROM games WHERE invite_code = $1", inviteCode,
	).Scan(&gameID, &status)
	if err != nil {
		return "", fmt.Errorf("invalid invite code: %w", err)
	}
	if status != "waiting" {
		return "", fmt.Errorf("game already started or finished")
	}

	// getGame acquires and releases the engine lock internally
	game, err := g.getGame(ctx, gameID)
	if err != nil {
		return "", err
	}

	// Now lock only the game object to mutate its Players map
	game.mu.Lock()
	game.Players[playerID] = &Player{
		Id:           playerID,
		Hints:        5,
		Score:        0,
		IsEliminated: false,
	}
	game.mu.Unlock()

	_, err = g.DBConn.ExecContext(ctx,
		"INSERT INTO game_players (game_id, user_id, hints_remaining) VALUES ($1, $2, $3)",
		gameID, playerID, 5,
	)
	if err != nil {
		return "", fmt.Errorf("could not add player to game: %w", err)
	}
	return gameID, nil
}

// StartGame transitions a game from waiting → active and begins the round loop.
func (g *GameEngine) StartGame(ctx context.Context, gameID string, hostID string) (string, error) {
	game, err := g.getGame(ctx, gameID)
	if err != nil {
		return "", err
	}

	// Read-only check — safe with RLock
	game.mu.RLock()
	isHost := game.HostId == hostID
	game.mu.RUnlock()

	if !isHost {
		return "", fmt.Errorf("only the host can start the game")
	}

	_, err = g.DBConn.ExecContext(ctx,
		"UPDATE games SET status = 'active' WHERE id = $1", gameID,
	)
	if err != nil {
		return "", fmt.Errorf("cannot update game status: %w", err)
	}

	game.mu.Lock()
	game.Status = "active"
	game.Start = time.Now()
	game.mu.Unlock()

	// Context lifecycle rule: HTTP request context dies when response is sent.
	// Long-running goroutines MUST use context.Background().
	go g.runRound(context.Background(), gameID)
	return gameID, nil
}

// selectLetter picks a random letter, ensuring no letter is used more than twice.
// Caller must hold game.mu.Lock() before calling this.
func (g *Game) selectLetter() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letter := make([]byte, 1)
	for {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		letter[0] = charset[n.Int64()]
		if g.LetterUsage[string(letter)] < 2 {
			g.LetterUsage[string(letter)]++
			break
		}
	}
	return string(letter)
}

// runRound drives the entire game loop — letter cycles, breaks, and scoring.
//
// Locking pattern here:
//   - We acquire game.mu.Lock() whenever we write to game fields
//   - We release before sleeping so other goroutines can read game state
//     (e.g. the WebSocket handler reading CurrentPhase for reconnects)
func (g *GameEngine) runRound(ctx context.Context, gameID string) {
	// Safe map read — use engine RLock since we're only reading, not modifying map
	g.mu.RLock()
	game := g.ActiveGames[gameID]
	g.mu.RUnlock()

	if game == nil {
		return
	}

	roundEnd := time.Now().Add(3 * time.Minute)

	game.mu.Lock()
	game.RoundEndsAt = roundEnd
	game.mu.Unlock()

	for time.Now().Before(roundEnd) {
		// Lock to read+write game fields atomically
		game.mu.Lock()
		letter := game.selectLetter() // selectLetter mutates LetterUsage — needs lock
		game.CurrentLetter = letter
		game.CurrentPhase = "playing"
		game.PhaseEndsAt = time.Now().Add(10 * time.Second)
		game.mu.Unlock()

		var roundID string
		err := g.DBConn.QueryRowContext(ctx,
			"INSERT INTO rounds (game_id, letter, started_at) VALUES ($1, $2, $3) RETURNING id",
			gameID, letter, time.Now(),
		).Scan(&roundID)
		if err != nil {
			log.Printf("failed to insert round: %v", err)
			continue
		}

		// Store roundID on game so reconnecting clients can get it via STATE
		game.mu.Lock()
		game.CurrentRoundID = roundID
		game.mu.Unlock()

		log.Printf("round inserted: %s", roundID)

		g.Hub.BroadcastMsg <- ws.BroadcastMessage{
			RoomId:  gameID,
			Message: []byte("ROUND:" + roundID),
		}
		g.Hub.BroadcastMsg <- ws.BroadcastMessage{
			RoomId:  gameID,
			Message: []byte("LETTER:" + letter),
		}

		// Release lock before sleeping — other goroutines need to read game state
		time.Sleep(10 * time.Second)

		game.mu.Lock()
		game.CurrentPhase = "break"
		game.PhaseEndsAt = time.Now().Add(5 * time.Second)
		game.mu.Unlock()

		g.Hub.BroadcastMsg <- ws.BroadcastMessage{
			RoomId:  gameID,
			Message: []byte("BREAK:5"),
		}

		// Broadcast current scores to all players after each break
		if err := g.broadcastScores(ctx, gameID); err != nil {
			log.Printf("failed to broadcast scores: %v", err)
		}

		time.Sleep(5 * time.Second)
	}

	g.eliminatePlayer(ctx, gameID)
}

// broadcastScores queries game_players and sends SCORES:{} to all clients.
// This implements the previously missing SCORES broadcast.
func (g *GameEngine) broadcastScores(ctx context.Context, gameID string) error {
	rows, err := g.DBConn.QueryContext(ctx,
		"SELECT user_id, score FROM game_players WHERE game_id = $1 AND is_eliminated = FALSE",
		gameID,
	)
	if err != nil {
		return fmt.Errorf("error fetching scores: %w", err)
	}
	defer rows.Close()

	scores := make(map[string]int)
	for rows.Next() {
		var userID string
		var score int
		if err := rows.Scan(&userID, &score); err != nil {
			return fmt.Errorf("error scanning scores: %w", err)
		}
		scores[userID] = score
	}

	// Build SCORES: payload manually to avoid importing encoding/json in a hot path
	
	scoresJSON, err := json.Marshal(scores)
	if err != nil {
		return fmt.Errorf("error marshalling scores: %w", err)
	}
	payload := append([]byte("SCORES:"), scoresJSON...)

	g.Hub.BroadcastMsg <- ws.BroadcastMessage{
		RoomId:  gameID,
		Message: []byte(payload),
	}
	return nil
}

// eliminatePlayer removes the lowest-scoring non-eliminated player.
// If more than 2 players remain, starts another round. Otherwise ends the game.
func (g *GameEngine) eliminatePlayer(ctx context.Context, gameID string) {
	g.mu.RLock()
	game := g.ActiveGames[gameID]
	g.mu.RUnlock()

	if game == nil {
		return
	}

	game.mu.Lock()

	var lowestID string
	lowestScore := math.MaxInt64
	for playerID, player := range game.Players {
		if !player.IsEliminated && player.Score < lowestScore {
			lowestScore = player.Score
			lowestID = playerID
		}
	}
	if lowestID != "" {
		game.Players[lowestID].IsEliminated = true
	}

	count := 0
	for _, p := range game.Players {
		if !p.IsEliminated {
			count++
		}
	}

	game.mu.Unlock()

	if count > 2 {
		go g.runRound(ctx, gameID)
	} else {
		game.mu.Lock()
		game.Status = "finished"
		game.mu.Unlock()

		g.Hub.BroadcastMsg <- ws.BroadcastMessage{
			RoomId:  gameID,
			Message: []byte("GAME:FINISHED"),
		}
	}
}

// GetGameState returns a safe snapshot of game state for reconnect purposes.
// Uses RLock since this is purely reading — multiple reconnects can happen simultaneously.
func (g *GameEngine) GetGameState(gameID string) (phase, letter, roundID string, timer, gameTime int) {
	g.mu.RLock()
	game := g.ActiveGames[gameID]
	g.mu.RUnlock()

	if game == nil {
		return
	}

	game.mu.RLock()
	defer game.mu.RUnlock()

	if game.CurrentPhase == "" {
		return
	}

	phase = game.CurrentPhase
	letter = game.CurrentLetter
	roundID = game.CurrentRoundID

	// Timer clamp: time.Until returns negative if phase already ended.
	// max(0, value) prevents the frontend from receiving -3 seconds remaining.
	rawTimer := int(time.Until(game.PhaseEndsAt).Seconds())
	if rawTimer < 0 {
		rawTimer = 0
	}
	timer = rawTimer

	rawGameTime := int(time.Until(game.RoundEndsAt).Seconds())
	if rawGameTime < 0 {
		rawGameTime = 0
	}
	gameTime = rawGameTime

	return
}
