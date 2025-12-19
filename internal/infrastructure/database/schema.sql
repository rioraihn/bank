CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       name varchar(50),
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE wallets (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         user_id UUID NOT NULL UNIQUE,
                         balance BIGINT NOT NULL DEFAULT 0,
                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),


                         CONSTRAINT wallets_balance_non_negative CHECK (balance >= 0),
                         CONSTRAINT wallets_user_fk FOREIGN KEY (user_id)
                             REFERENCES users(id)
                             ON DELETE CASCADE
);

CREATE TABLE transactions (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              wallet_id UUID NOT NULL,
                              transaction_type VARCHAR(20) NOT NULL,
                              amount BIGINT NOT NULL,
                              created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),


                              CONSTRAINT transactions_type_valid CHECK (
                                  transaction_type IN ('WITHDRAWAL', 'DEPOSIT')
                                  ),
                              CONSTRAINT transactions_amount_positive CHECK (amount > 0),
                              CONSTRAINT transactions_wallet_fk FOREIGN KEY (wallet_id)
                                  REFERENCES wallets(id)
                                  ON DELETE CASCADE
);

CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_transactions_wallet_id ON transactions(wallet_id);

-- 1. Create user
INSERT INTO users (id, name) VALUES ('cfa3b5c8-258a-4d9a-9258-d0ab849ef82d', 'rio');
INSERT INTO users (id, name) VALUES ('cfa3b5c8-258a-4d9a-9258-d0ab849ef82f', 'raihan');


-- 2. Create wallet
INSERT INTO wallets (id, user_id, balance) VALUES (
                                                      '11111111-1111-1111-1111-111111111111',
                                                      'cfa3b5c8-258a-4d9a-9258-d0ab849ef82d',
                                                      100000
                                                  );

INSERT INTO wallets (id, user_id, balance) VALUES (
                                                      '11111111-1111-1111-1111-111111111112',
                                                      'cfa3b5c8-258a-4d9a-9258-d0ab849ef82f',
                                                      500000
                                                  );
