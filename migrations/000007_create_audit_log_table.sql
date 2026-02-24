CREATE TABLE audit_log(
    id UUID PRIMARY KEY DEFAULT gen_random_UUID(),
    submission_id UUID REFERENCES submissions(id),
    reviewed_by UUID REFERENCES users(id),
    decision TEXT NOT NULL DEFAULT 'pending',
    reason_for_decision TEXT,
    decided_at TIMESTAMP WITH TIME ZONE DEFAULT now()

)