package game

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"

	ws "github.com/jerryjuche/alpha-game/internal/websocket"
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
}

func NewGameEngine(db *sqlx.DB, h *ws.Hub) *GameEngine {
	return &GameEngine{
		DBConn:      db,
		Hub:         h,
		ActiveGames: make(map[string]*Game),
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

func (g *GameEngine) CreateGame(ctx context.Context, hostID string) (string, string, error) {
	var gameID string
	var inviteCode string

	inviteCode = generateInviteCode()

	err := g.DBConn.QueryRowContext(ctx, "INSERT INTO games (invite_code, status, created_by) VALUES ($1, $2, $3) RETURNING id", inviteCode, "waiting", hostID).Scan(&gameID)
	if err != nil {
		return "", "", fmt.Errorf("cannot create game: %w", err)

	}

	ActiveGame := &Game{
		GameId:        gameID,
		LetterUsage:   make(map[string]int),
		Start:         time.Now(),
		End:           time.Now(),
		Status:        "waiting",
		CurrentLetter: "",
		HostId:        hostID,
	}

	g.ActiveGames[gameID] = ActiveGame
	return gameID, inviteCode, nil

}

func (g *GameEngine) JoinGame(ctx context.Context, playerID string, inviteCode string) (string, error) {
	var gameID string
	var status string

	err := g.DBConn.QueryRowContext(ctx, "SELECT id, status FROM games WHERE invite_code = $1", inviteCode).Scan(&gameID, &status)
	if err != nil {
		return "", fmt.Errorf("Invalid Invite code, %w", err)
	}
	if status != "waiting" {
		return "", fmt.Errorf("Game started or finished, %w", err)
	}

	game := g.ActiveGames[gameID]

	game.Players[playerID] = &Player{
		Id:           playerID,
		Hints:        5,
		Score:        0,
		IsEliminated: false,
	}

	_, err = g.DBConn.ExecContext(ctx, "INSERT INTO game_players (game_id, user_id, hints_remaining) VALUES ($1, $2, $3)", gameID, playerID, 5)
	if err != nil {
		return "", fmt.Errorf("Invalid user, %w", err)
	}

	return gameID, nil

}

func (g *GameEngine) StartGame(ctx context.Context, gameID string, hostID string, status string) (string, error) {

	_, err := g.DBConn.ExecContext(ctx, "UPDATE games SET is_status = 'active' WHERE id = $1", gameID)
	if err != nil {
		return "", fmt.Errorf("Cannou update table, %w", err)
	}

	game := g.ActiveGames[gameID]
	if game.HostId != hostID {
		return "only the host can start the game", nil
	}

	game.Status = "active"
	game.Start = time.Now()

	go g.runRound(gameID)

	return gameID, nil

}

func (g *Game) selectLetter() string {

	// Generate random letter
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

func (g *GameEngine) runRound(gameID string) {
	game := g.ActiveGames[gameID]

	// 5 mins timer
	roundEnd := time.Now().Add(5 * time.Minute)

	for time.Now().Before(roundEnd) {
		letter := game.selectLetter()

		g.Hub.BroadcastMsg <- ws.BroadcastMessage{
			RoomId:  gameID,
			Message: []byte("LETTER:" + letter),
		}

		time.Sleep(8 * time.Second)

	}

}

func (g *GameEngine) eliminatePlayer(gameID string) {

	game := g.ActiveGames[gameID]

	var lowestID string
	lowestScore := math.MaxInt64

	for playerID, player := range game.Players {
		if player.Score < lowestScore {
			lowestScore = player.Score // update lowest score!
			lowestID = playerID
		}
	}
	game.Players[lowestID].IsEliminated = true // eliminate AFTER loop

	count := 0
	for _, p := range game.Players {
		if !p.IsEliminated {
			count++
		}
		if count > 2 {
			go g.runRound(gameID)
		} else {
			game.Status = "finished"
		}
	}

}
