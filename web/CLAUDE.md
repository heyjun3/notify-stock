# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

```bash
# Development
bun run dev              # Start development server with HMR at http://localhost:5173

# Building and Type Checking
bun run build           # Build for production using React Router
bun run typecheck       # Generate React Router types and run TypeScript check

# Code Quality
bun run fmt             # Format code using Biome
bun run lint            # Lint and auto-fix using Biome

# GraphQL Code Generation
bun run codegen         # Generate TypeScript types from GraphQL schema and format

# Production
bun run start           # Start production server from build output
```

## Architecture Overview

This is a React Router v7 application with the following key architectural components:

### Frontend Stack
- **React Router v7**: Full-stack React framework with SSR
- **Apollo Client**: GraphQL client for data fetching
- **TailwindCSS**: Utility-first CSS framework
- **Recharts**: React charting library for stock price visualization
- **Biome**: Code formatting and linting

### Project Structure
- `app/` - Main application code following React Router conventions
- `app/gen/graphql.ts` - Auto-generated GraphQL types and hooks
- `app/dashboard/` - Dashboard-specific components (stock cards, charts, pagination)
- `app/routes/` - Route components
- `build/` - Production build output

### Data Flow
- GraphQL schema served from `http://localhost:8080/query` (backend API)
- Apollo Client configured in `app/root.tsx` with backend URL from `VITE_BACKEND_URL`
- GraphQL queries defined in `.gql` files (e.g., `app/dashboard/getSymbol.gql`)
- Auto-generated hooks using GraphQL Code Generator

### Key Configuration
- **Backend URL**: Configured via `VITE_BACKEND_URL` environment variable
- **GraphQL Codegen**: Generates types from `http://localhost:8080/query` schema
- **Path Aliases**: `~/*` maps to `./app/*` in TypeScript config
- **Apollo CJS Interop**: Required for Apollo Client compatibility with Vite

### Application Features
- Stock price dashboard with search and pagination
- Interactive stock price charts with period selection (1M, 6M, 1Y, 5Y)
- Real-time stock data visualization using Recharts
- Responsive design with dark mode support

## Development Workflow

1. Ensure backend GraphQL API is running at `http://localhost:8080/query`
2. Run `bun run codegen` when GraphQL schema changes
3. Use `bun run typecheck` before committing to catch type errors
4. Run `bun run fmt && bun run lint` for code quality checks
