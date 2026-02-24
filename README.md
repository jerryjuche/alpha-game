# AlphaBlitz ⚡

> A fast-paced, real-time multiplayer word game built in Go.

AlphaBlitz challenges players to fill in words across five categories — **Name, Animal, Place, Thing, and Food** — all starting with a randomly selected letter. With an 8-second timer per letter and a 5-minute round clock, every second counts.

---

## 🎮 How It Works

A random letter is selected every 8 seconds. Players race to fill in a valid word for each of the five categories before time runs out. Each correct answer earns **10 points**. After each 5-minute round, the player with the lowest score is eliminated. The last two players standing face off in a final showdown.

**Example — Letter "A":**

| Category | Answer     |
|----------|------------|
| Name     | Angela     |
| Animal   | Antelope   |
| Place    | Afghanistan|
| Thing    | Aeroplane  |
| Food     | Amala      |

---

## ✨ Features

- Real-time multiplayer with WebSocket connections (up to 5 players per room)
- Shareable invite codes for room joining
- 8-second letter rotation with 5-minute round timer
- Letters can only repeat twice per game
- Automated answer submission on timer expiry
- Audit system — all answers reviewed against a word database
- Partial credit (5pts) for misspelled but valid words
- 5 hints per player — auto-fills one empty field
- Admin dashboard for word database management
- Bulk word imports via Excel files (.xlsx)
- Player profiles with full game statistics
- JWT-based authentication
- Light and dark theme toggle
- Elimination system — lowest scorer removed each round

---

## 🏗️ Architecture

```
alphablitz/
├── cmd/server/          # Entry point
├── internal/
│   ├── auth/            # Registration, login, JWT middleware
│   ├── game/            # Game engine, room management, scoring
│   ├── websocket/       # Hub, client connections, real-time broadcasting
│   ├── audit/           # Answer review and point assignment
│   ├── word/            # Word database, Excel importer
│   ├── user/            # Player profiles and stats
│   └── repository/      # Database access layer
├── pkg/
│   ├── jwt/             # Token utilities
│   ├── validator/       # Input validation
│   └── response/        # Standardized API responses
├── migrations/          # SQL migration files
├── config/              # Environment configuration
└── web/                 # Frontend
```

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22 |
| Router | chi |
| Database | PostgreSQL |
| Query Layer | sqlx + pgx |
| WebSockets | gorilla/websocket |
| Auth | JWT (golang-jwt/jwt) |
| Password Hashing | bcrypt |
| Excel Import | excelize |
| Config | godotenv |
| Frontend | HTML + CSS + Vanilla JS |

---

## 🗄️ Database Schema

```
users           — accounts, profiles, theme preferences
games           — rooms, invite codes, game status
game_players    — player state per game (score, hints, elimination)
rounds          — round history with letter and timestamps
submissions     — player answers per round
word_database   — approved word dictionary by category and letter
audit_log       — auditor decisions on all submissions
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL 14+

### Setup

```bash
# Clone the repository
git clone https://github.com/jerryjuche/alpha-game.git
cd alpha-game

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your database credentials and JWT secret

# Run migrations
psql -d alphablitz -f migrations/000001_create_users_table.sql
psql -d alphablitz -f migrations/000002_create_games_table.sql
psql -d alphablitz -f migrations/000003_create_game_players_table.sql
psql -d alphablitz -f migrations/000004_create_rounds_table.sql
psql -d alphablitz -f migrations/000005_create_submissions_table.sql
psql -d alphablitz -f migrations/000006_create_word_database_table.sql
psql -d alphablitz -f migrations/000007_create_audit_log_table.sql

# Run the server
go run cmd/server/main.go
```

---

## 📡 API Endpoints

### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Create a new account |
| POST | `/auth/login` | Login and receive JWT token |

### Game (Protected — requires JWT)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/game/create` | Create a new game room |
| POST | `/game/join` | Join a game via invite code |
| POST | `/game/start` | Host starts the game |

### WebSocket
| Endpoint | Description |
|----------|-------------|
| `ws://host/ws/:gameID` | Connect to live game room |

---

## 🎯 Game Rules

- Up to **5 players** per game
- Each letter is played for **8 seconds**
- Each round lasts **5 minutes**
- Each letter can only be selected **twice** per game
- **10 points** for a correct answer
- **5 points** for a misspelled but valid word (auditor discretion)
- After each round, the **lowest scorer is eliminated**
- Game ends when **2 players remain** for the final round
- Each player has **5 hints** per game

---

## 🔐 Environment Variables

```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=alphablitz
JWT_SECRET=your_super_secret_key
ENV=development
```

---

## 📊 Project Status

- [x] Project architecture & database design
- [x] Authentication system (register, login, JWT middleware)
- [x] Word database with Excel bulk import
- [x] Game engine (create, join, start, timer, elimination)
- [x] Real-time WebSocket infrastructure
- [ ] Submission handler
- [ ] Audit system
- [ ] Player profiles & statistics
- [ ] Frontend
- [ ] Admin dashboard

---

## 👨‍💻 Author

Built by [@jerryjuche](https://github.com/jerryjuche) — learning Go by building real things.

---

## 📄 License

MIT