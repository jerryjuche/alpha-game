CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    round_id UUID REFERENCES rounds(id),
    submitted_by UUID REFERENCES users(id),
    category TEXT NOT NULL,
    word_submitted TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    points_awarded INTEGER DEFAULT 0,
    CONSTRAINT unique_player_round_category UNIQUE (round_id, submitted_by, category)
);