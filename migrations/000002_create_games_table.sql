CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_UUID(),
    invite_code TEXT NOT NULL,
    status TEXT DEFAULT 'status',
    active_letter TEXT NOT NULL,
    current_round INTEGER DEFAULT 0 ,
    created_by UUID REFERENCES users(id) ,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);