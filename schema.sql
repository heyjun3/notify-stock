CREATE TABLE IF NOT EXISTS
    stocks (
        symbol TEXT,
        TIMESTAMP TIMESTAMP,
        open DECIMAL NOT NULL,
        CLOSE DECIMAL NOT NULL,
        high DECIMAL NOT NULL,
        low DECIMAL NOT NULL,
        PRIMARY KEY (symbol, TIMESTAMP)
    )
PARTITION BY
    LIST (symbol);

CREATE TABLE IF NOT EXISTS
    stocks_n225 PARTITION OF stocks FOR
VALUES
    IN ('N225');

CREATE TABLE IF NOT EXISTS
    stocks_sp500 PARTITION OF stocks FOR
VALUES
    IN ('S&P500');

CREATE TABLE IF NOT EXISTS
    notifications (
        id UUID PRIMARY KEY,
        symbol TEXT,
        Email TEXT,
        Hour TIME
    );
