CREATE TABLE rounds (
    id UUID PRIMARY KEY DEFAULT gen_random_UUID(),
    game_id UUID REFERENCES games(id),
    letter TEXT NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    ended_at TIMESTAMP WITH TIME ZONE 
)