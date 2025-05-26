CREATE TABLE IF NOT EXISTS
    stocks (
        symbol TEXT,
        TIMESTAMP TIMESTAMP,
        open DECIMAL NOT NULL,
        CLOSE DECIMAL NOT NULL,
        high DECIMAL NOT NULL,
        low DECIMAL NOT NULL,
        PRIMARY KEY (symbol, TIMESTAMP)
    );

CREATE TABLE IF NOT EXISTS 
    symbols (
        symbol TEXT PRIMARY KEY,
        short_name TEXT NOT NULL,
        long_name TEXT NOT NULL,
        market_price DECIMAL NOT NULL,
        previous_close DECIMAL NOT NULL,
        volume INTEGER,
        market_cap INTEGER
    )
;
ALTER TABLE symbols ADD COLUMN currency TEXT;

CREATE TABLE IF NOT EXISTS
    notifications (
        id UUID PRIMARY KEY,
        symbol TEXT,
        Email TEXT,
        Hour TIME
    );

CREATE TABLE IF NOT EXISTS 
    notification_targets (
        id UUID PRIMARY KEY,
        symbol TEXT,
        notification_id UUID,
        FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE
    );

ALTER TABLE notification_targets ADD CONSTRAINT fk_symbol
    FOREIGN KEY (symbol)
    REFERENCES symbols(symbol)
    ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS
    sessions (
        id TEXT PRIMARY KEY,
        state TEXT NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        expires_at TIMESTAMP NOT NULL
    );
