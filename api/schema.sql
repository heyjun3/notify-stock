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
    notifications (
        id UUID PRIMARY KEY,
        symbol TEXT,
        Email TEXT,
        Hour TIME
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
