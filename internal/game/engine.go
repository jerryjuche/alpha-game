package game

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"

	ws "github.com/jerryjuche/alpha-game/internal/websocket"
	"github.com/jerryjuche/alpha-game/internal/word"
	"github.com/jmoiron/sqlx"
)

type Game struct {
	Players       map[string]*Player
	GameId        string
	LetterUsage   map[string]int
	Start         time.Time
	End           time.Time
	Status        string
	CurrentLetter string
	HostId        string
	PhaseEndsAt   time.Time
	RoundEndsAt   time.Time
	CurrentPhase  string
}

type Player struct {
	Id           string
	Hints        int
	Score        int
	IsEliminated bool
}

type GameEngine struct {
	DBConn      *sqlx.DB
	Hub         *ws.Hub
	ActiveGames map[string]*Game
	WordService *word.WordService
}

func NewGameEngine(db *sqlx.DB, h *ws.Hub, w *word.WordService) *GameEngine {
	return &GameEngine{
		DBConn:      db,
		Hub:         h,
		ActiveGames: make(map[string]*Game),
		WordService: w,
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

func (g *GameEngine) getGame(ctx context.Context, gameID string) (*Game, error) {
	game := g.ActiveGames[gameID]
	if game == nil {
		if err := g.loadGameFromDB(ctx, gameID); err != nil {
			return nil, fmt.Errorf("game not found: %w", err)
		}
		game = g.ActiveGames[gameID]
	}
	return game, nil
}

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

	return gameID, inviteCode, nil
}

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

	game, err := g.getGame(ctx, gameID)
	if err != nil {
		return "", err
	}

	game.Players[playerID] = &Player{
		Id:           playerID,
		Hints:        5,
		Score:        0,
		IsEliminated: false,
	}

	_, err = g.DBConn.ExecContext(ctx,
		"INSERT INTO game_players (game_id, user_id, hints_remaining) VALUES ($1, $2, $3)",
		gameID, playerID, 5,
	)
	if err != nil {
		return "", fmt.Errorf("could not add player to game: %w", err)
	}

	return gameID, nil
}

func (g *GameEngine) StartGame(ctx context.Context, gameID string, hostID string) (string, error) {
	game, err := g.getGame(ctx, gameID)
	if err != nil {
		return "", err
	}

	if game.HostId != hostID {
		return "", fmt.Errorf("only the host can start the game")
	}

	_, err = g.DBConn.ExecContext(ctx,
		"UPDATE games SET status = 'active' WHERE id = $1", gameID,
	)
	if err != nil {
		return "", fmt.Errorf("cannot update game status: %w", err)
	}

	game.Status = "active"
	game.Start = time.Now()

	go g.runRound(ctx, gameID)

	return gameID, nil
}

func (g *Game) selectLetter() string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
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

func (g *GameEngine) runRound(ctx context.Context, gameID string) {
	game := g.ActiveGames[gameID]
	if game == nil {
		return
	}

	roundEnd := time.Now().Add(3 * time.Minute)
	game.RoundEndsAt = roundEnd

	for time.Now().Before(roundEnd) {

		letter := game.selectLetter()
		game.CurrentLetter = letter
		game.CurrentPhase = "playing"
		game.PhaseEndsAt = time.Now().Add(10 * time.Second)

		var roundID string

		err := g.DBConn.QueryRowContext(ctx, "INSERT INTO rounds (game_id, letter, started_at) VALUES ($1, $2, $3) RETURNING id", gameID, letter, time.Now()).Scan(&roundID)
		if err != nil {
			fmt.Errorf("error fetching data", err)
		}

		g.Hub.BroadcastMsg <- ws.BroadcastMessage{
			RoomId:  gameID,
			Message: []byte("LETTER:" + letter),
		}

		time.Sleep(10 * time.Second)

		game.CurrentPhase = "break"
		game.PhaseEndsAt = time.Now().Add(5 * time.Second)

		g.Hub.BroadcastMsg <- ws.BroadcastMessage{
			RoomId:  gameID,
			Message: []byte("BREAK:5"),
		}

		time.Sleep(5 * time.Second)
	}

	g.eliminatePlayer(ctx, gameID)
}

func (g *GameEngine) eliminatePlayer(ctx context.Context, gameID string) {
	game := g.ActiveGames[gameID]
	if game == nil {
		return
	}

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

	if count > 2 {
		go g.runRound(ctx, gameID)
	} else {
		game.Status = "finished"
		g.Hub.BroadcastMsg <- ws.BroadcastMessage{
			RoomId:  gameID,
			Message: []byte("GAME:FINISHED"),
		}
	}
}
