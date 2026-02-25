CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_UUID(),
    round_id UUID REFERENCES rounds(id),
    submitted_by UUID REFERENCES users(id),
    category TEXT NOT NULL,
    word_submitted TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    points_awarded INTEGER
)