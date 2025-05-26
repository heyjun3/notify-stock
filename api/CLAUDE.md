# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Stock notification API built in Go that fetches stock market data, stores it in PostgreSQL, and sends email notifications. Features a GraphQL API, CLI commands, and Docker deployment.

## Common Commands

### Development Setup
```bash
make db-setup          # Start PostgreSQL container and create database schemas
make build             # Build Docker image
make gqlgen           # Generate GraphQL code from schema
```

### Daily Operations
```bash
make notify           # Send email notifications for N225 and S&P 500
make update           # Fetch last 7 days of stock data
make update-all       # Fetch 5 years of historical data
make db-connect       # Connect to development database
```

### Testing
```bash
go test ./internal/... # Run all tests in internal package
go test -v ./internal/stock_test.go # Run specific test file
```

### Development Server
```bash
go run cmd/main.go server  # Start GraphQL server on localhost:8080
```

## Architecture

### CLI Structure
The project uses Cobra CLI with commands in `cmd/` directory:
- `server` - GraphQL API server
- `notify` - Email notification system
- `stock update` - Stock data fetching
- `email`, `fetch`, `version`, `yaml` - Utility commands

### Core Components
- **internal/** - Business logic with dependency injection via Wire
- **graph/** - GraphQL schema, resolvers, and generated code
- **cmd/** - CLI command implementations

### Data Flow
1. `stock update` fetches data from Yahoo Finance → PostgreSQL
2. `notify` queries database → generates email summaries → MailerSend
3. GraphQL API serves data with OAuth authentication

### Database Schema
- `stocks` - OHLC price data with timestamps
- `symbols` - Stock metadata (name, current price, volume)
- `notifications` - User notification preferences
- `notification_targets` - User watchlists

### Configuration
- `config.yaml` - Supported stock symbols (^N225, ^GSPC, ^DJI, ^IXIC, ^XDN)
- Environment variables for email config, database DSN, OAuth settings
- `compose.yaml` - PostgreSQL service on port 5555

### Key Files
- `schema.sql` - Database migrations
- `graph/schema.graphqls` - GraphQL schema definition
- `internal/wire.go` - Dependency injection setup
- `gqlgen.yml` - GraphQL code generation config

## Development Notes

### GraphQL Development
Run `make gqlgen` after modifying `graph/schema.graphqls` to regenerate resolvers and types.

### Authentication
All GraphQL mutations and some queries require OAuth authentication. Session middleware handles auth state.

### Stock Symbols
Add new symbols to `config.yaml` and update database via `symbols` table. The system supports any Yahoo Finance symbol.

### Email Notifications
Use MailerSend API. Configure FROM/TO emails and MAIL_TOKEN environment variables.