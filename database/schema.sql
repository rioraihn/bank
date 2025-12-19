-- Banking Service Database Schema
-- PostgreSQL Database Setup
-- Host: localhost:5433
-- User: rio
-- Password: rio
-- Database: postgres

-- Drop existing tables if they exist (for fresh setup)
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS wallets CASCADE;

-- Create wallets table
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT wallets_balance_positive CHECK (balance >= 0)
);

-- Create transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    failure_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT transactions_type_valid CHECK (transaction_type IN ('WITHDRAWAL', 'DEPOSIT')),
    CONSTRAINT transactions_status_valid CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED')),
    CONSTRAINT transactions_amount_positive CHECK (amount > 0),

    -- Foreign Key
    FOREIGN KEY (wallet_id) REFERENCES wallets(user_id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_transactions_wallet_id ON transactions(wallet_id);
CREATE INDEX idx_transactions_type ON transactions(transaction_type);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_wallets_updated_at
    BEFORE UPDATE ON wallets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data for testing
INSERT INTO wallets (user_id, balance) VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 100000), -- $1000.00
    ('550e8400-e29b-41d4-a716-446655440001', 50000);   -- $500.00

-- Verify the setup
SELECT 'Wallets table created' as status;
SELECT COUNT(*) as wallet_count FROM wallets;