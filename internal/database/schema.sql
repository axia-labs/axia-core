-- Schema for Axia Trust Graph Database

-- Claims table stores all axiomatic claims
CREATE TABLE claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issuer VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    axiom_text TEXT NOT NULL,
    confidence DECIMAL(4,3) NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    proof_type VARCHAR(100) NOT NULL,
    proof_value TEXT NOT NULL,
    proof_created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Tags for claims
CREATE TABLE claim_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Trust graph edges
CREATE TABLE trust_edges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_node UUID NOT NULL REFERENCES claims(id),
    to_node UUID NOT NULL REFERENCES claims(id),
    weight DECIMAL(4,3) NOT NULL CHECK (weight >= 0 AND weight <= 1),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Twitter reports
CREATE TABLE twitter_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tweet_id VARCHAR(255) NOT NULL UNIQUE,
    author_id VARCHAR(255) NOT NULL,
    subject_handle VARCHAR(255) NOT NULL,
    action VARCHAR(100) NOT NULL,
    project VARCHAR(100) NOT NULL,
    amount VARCHAR(50),
    claim_id UUID NOT NULL REFERENCES claims(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- IPFS records table
CREATE TABLE ipfs_records (
    id UUID PRIMARY KEY,
    ipfs_id VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes
CREATE INDEX idx_claims_issuer ON claims(issuer);
CREATE INDEX idx_claims_subject ON claims(subject);
CREATE INDEX idx_claim_tags_claim_id ON claim_tags(claim_id);
CREATE INDEX idx_claim_tags_tag ON claim_tags(tag);
CREATE INDEX idx_trust_edges_from_node ON trust_edges(from_node);
CREATE INDEX idx_trust_edges_to_node ON trust_edges(to_node);
CREATE INDEX idx_twitter_reports_tweet_id ON twitter_reports(tweet_id);
CREATE INDEX idx_ipfs_records_type ON ipfs_records(type);
CREATE INDEX idx_ipfs_records_created_at ON ipfs_records(created_at); 