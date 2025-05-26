# Stock Notification API

A Go-based stock market monitoring system that fetches real-time stock data, stores historical prices, and sends automated email notifications with market summaries.

## Features

- **Stock Data Fetching** - Retrieves market data from Yahoo Finance for major indices
- **Email Notifications** - Automated daily market summaries via MailerSend
- **GraphQL API** - Query stock data and manage notification preferences
- **OAuth Authentication** - Secure user sessions and API access
- **Docker Deployment** - Complete containerized environment
- **CLI Interface** - Command-line tools for operations and maintenance

## Supported Markets

- Nikkei 225 (^N225)
- S&P 500 (^GSPC) 
- Dow Jones (^DJI)
- NASDAQ (^IXIC)
- Custom indices

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for development)

### Setup

1. **Start the database**
   ```bash
   make db-setup
   ```

2. **Configure environment variables**
   ```bash
   export DBDSN="postgres://user:password@localhost:5555/stockdb?sslmode=disable"
   export MAIL_TOKEN="your-mailersend-token"
   export FROM="sender@example.com"
   export TO="recipient@example.com"
   export OAUTH_CLIENT_ID="your-oauth-client-id"
   export OAUTH_CLIENT_SECRET="your-oauth-secret"
   export OAUTH_REDIRECT_URL="http://localhost:8080/auth/callback"
   ```

3. **Build the application**
   ```bash
   make build
   ```

### Usage

#### Start the GraphQL API Server
```bash
go run cmd/main.go server
```
Access GraphQL playground at http://localhost:8080

#### Fetch Stock Data
```bash
# Fetch last 7 days
go run cmd/main.go stock update

# Fetch all historical data (5 years)
go run cmd/main.go stock update -a
```

#### Send Email Notifications
```bash
# Send notifications for N225 and S&P 500
go run cmd/main.go notify -s "^N225,^GSPC"
```

## API Documentation

### GraphQL Queries

```graphql
# Get all symbols
query {
  symbols {
    symbol
    name
    currentPrice
    volume
  }
}

# Get specific symbol with price history
query {
  symbol(symbol: "^GSPC") {
    symbol
    name
    stocks(limit: 10) {
      date
      open
      high
      low
      close
      volume
    }
  }
}

# Get user notifications (requires auth)
query {
  notifications {
    id
    name
    targets {
      symbol
    }
  }
}
```

### GraphQL Mutations

```graphql
# Create notification (requires auth)
mutation {
  createNotification(input: {
    name: "My Portfolio"
    symbols: ["^GSPC", "^N225"]
  }) {
    id
    name
  }
}
```

## Development

### Database Schema

The application uses PostgreSQL with the following main tables:

- `symbols` - Stock metadata and current prices
- `stocks` - Historical OHLC price data
- `notifications` - User notification preferences  
- `notification_targets` - Symbol watchlists per notification

### Adding New Stock Symbols

1. Add symbol to `config.yaml`
2. Insert symbol metadata into `symbols` table
3. Run `stock update` to fetch historical data

### GraphQL Development

After modifying `graph/schema.graphqls`:
```bash
make gqlgen
```

### Testing

```bash
go test ./internal/...
```

## Docker Deployment

### Using Docker Compose

```bash
docker-compose up -d
```

This starts:
- PostgreSQL database on port 5555
- Application server on port 8080

### Manual Docker Build

```bash
make build
docker run -p 8080:8080 notify-stock-api
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `DBDSN` | PostgreSQL connection string | Yes |
| `MAIL_TOKEN` | MailerSend API token | Yes |
| `FROM` | Sender email address | Yes |
| `TO` | Recipient email address | Yes |
| `OAUTH_CLIENT_ID` | OAuth client ID | Yes |
| `OAUTH_CLIENT_SECRET` | OAuth client secret | Yes |
| `OAUTH_REDIRECT_URL` | OAuth redirect URL | Yes |

### Stock Symbols Configuration

Edit `config.yaml` to modify supported stock symbols:

```yaml
symbols:
  - ^N225    # Nikkei 225
  - ^GSPC    # S&P 500
  - ^DJI     # Dow Jones
  - ^IXIC    # NASDAQ
  - ^XDN     # Custom index
```

## Architecture

The project follows clean architecture principles:

- **cmd/** - CLI commands and entry points
- **internal/** - Business logic and domain models
- **graph/** - GraphQL schema and resolvers
- **schema.sql** - Database migrations

Key design patterns:
- Dependency injection via Wire
- Repository pattern for data access
- Command pattern for CLI operations
- Session-based authentication middleware

## License

MIT License