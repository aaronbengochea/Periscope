# Periscope Go Backend

The Go backend handles general API operations, user management, portfolio tracking, and data aggregation for the Periscope options trading platform.

## Architecture

This backend follows the Standard Go Project Layout:

```
backend-go/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/                    # Private application code
│   ├── api/
│   │   ├── handlers/           # HTTP request handlers
│   │   ├── middleware/         # HTTP middleware (CORS, logging, auth)
│   │   └── router.go           # Route definitions
│   ├── services/               # Business logic
│   └── models/                 # Data structures
├── pkg/                        # Public libraries (reusable)
│   ├── database/               # Database connection
│   ├── massive/                # Massive API client
│   └── errors/                 # Error types
├── config/                     # Configuration management
│   └── config.go
├── Makefile                    # Build automation
└── go.mod                      # Go module definition
```

## Tech Stack

- **Framework**: Gin (HTTP web framework)
- **Database**: Supabase (PostgreSQL) via pgx
- **Config**: Viper (environment variables)
- **External APIs**: Massive.com (options data)

## Prerequisites

- Go 1.23 or higher
- Access to Supabase database
- Massive API key

## Setup

1. **Environment Variables**

   The backend reads configuration from the `.env` file in the project root:

   ```bash
   # From project root, not backend-go directory
   cp .env.example .env
   # Edit .env with your credentials
   ```

2. **Install Dependencies**

   ```bash
   go mod download
   ```

3. **Build**

   ```bash
   make build
   ```

4. **Run**

   ```bash
   make run
   ```

   The server will start on `http://localhost:8080`

## Available Commands

```bash
make build          # Build the application
make run            # Run the application
make test           # Run tests
make lint           # Run linter
make fmt            # Format code
make clean          # Clean build artifacts
```

## API Endpoints

### Health Check
```
GET /health
```

Returns server health status and database connection status.

### Options API (v1)
```
GET /api/v1/options/:ticker
```

Fetch options chain for a given ticker (to be implemented).

## Development

### Hot Reload

For development with hot reload, install `air`:

```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Testing

Run tests:
```bash
make test
```

Run with coverage:
```bash
make test-coverage
```

### Code Quality

Format code:
```bash
make fmt
```

Run linter:
```bash
make lint
```

## Project Structure Explanation

- **cmd/**: Application entry points (main packages)
- **internal/**: Private application code (cannot be imported by other projects)
  - `api/`: HTTP layer (routing, handlers, middleware)
  - `services/`: Business logic
  - `models/`: Data structures
- **pkg/**: Public libraries that can be imported by other projects
- **config/**: Configuration loading and validation

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `MASSIVE_API_KEY` | Massive.com API key | Yes |
| `MASSIVE_BASE_URL` | Massive API base URL | No (default: https://api.massive.com/v3) |
| `SUPABASE_URL` | Supabase project URL | Yes |
| `SUPABASE_ANON_KEY` | Supabase anonymous key | Yes |
| `SUPABASE_SERVICE_KEY` | Supabase service role key | Yes |
| `PORT` | Server port | No (default: 8080) |
| `GIN_MODE` | Gin mode (debug/release) | No (default: debug) |

## Next Steps

1. Implement Massive API client (`pkg/massive/client.go`)
2. Create data models (`internal/models/options.go`)
3. Implement options handlers (`internal/api/handlers/options.go`)
4. Add database queries with SQLC
5. Implement authentication middleware
6. Add comprehensive tests

## Separation from Rust Backend

The Rust execution backend (in `../src/`) handles:
- High-frequency order execution
- Low-latency market data processing
- Critical path optimization

This Go backend handles:
- REST API for frontend
- User management
- Portfolio tracking
- Data aggregation
- Business logic

Both backends can communicate via gRPC or message queues when needed.
