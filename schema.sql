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

CREATE TABLE
    stocks_n255 PARTITION OF stocks FOR
VALUES
    IN ('N255');

CREATE TABLE
    stocks_sp500 PARTITION OF stocks FOR
VALUES
    IN ('S&P500');