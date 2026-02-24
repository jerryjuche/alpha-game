CREATE TABLE word_database (
    id UUID PRIMARY KEY DEFAULT gen_random_UUID(),
    word TEXT NOT NULL,
    category TEXT NOT NULL,
    starts_with TEXT NOT NULL,
    added_by UUID REFERENCES users(id), 
    is_approved BOOLEAN DEFAULT FALSE,
    time_added TIMESTAMP WITH TIME ZONE DEFAULT now()

)