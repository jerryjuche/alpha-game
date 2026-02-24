CREATE TABLE game_players (
    id UUID PRIMARY KEY DEFAULT gen_random_UUID(),
    user_id UUID REFERENCES users(id),
    game_id UUID REFERENCES games(id),
    score INTEGER DEFAULT 0,
    hints_remaining INTEGER DEFAULT 5,
    is_eliminated BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);