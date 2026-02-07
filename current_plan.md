# Periscope - Go Backend & Frontend Foundation Plan

## Context

Periscope is transitioning from a single Rust-based analytics engine to a **dual-backend architecture**:
- **Go Backend**: Handle general API operations, user management, portfolio tracking, data aggregation
- **Rust Execution Backend**: High-performance order execution, real-time Greeks, critical path optimization

This plan focuses on establishing the **foundational infrastructure** for the Go backend and frontend, with 5 core objectives:

1. **Go Backend Core Structure** - Modern, scalable architecture
2. **React TypeScript Frontend** - Testing ground for data flow
3. **Massive API Integration** - Market data pipeline to Go backend
4. **Request/Response Objects** - Structured API contracts
5. **CI/CD Pipeline** - Automated testing and deployment

---

## Finalized Architecture Decisions ✅

Based on discussion, the following stack has been chosen:

1. **Go Framework**: ✅ **Gin** - Fast, excellent docs, production-ready
2. **Database**: ✅ **Supabase (PostgreSQL)** - Managed database with built-in auth, real-time, and REST API
3. **Database Layer**: ✅ **SQLC** - Type-safe SQL, compile-time validation, zero overhead
4. **Frontend**: ✅ **Next.js 15** - Full-featured framework with SSR, App Router, production-ready
5. **Development Approach**: ✅ **Speed First** - Supabase for managed DB, iterate quickly, scale as needed

### State Management for Next.js
- **Server State**: TanStack Query (React Query) for API data caching
- **Client State**: Zustand for UI state (minimal, only what's needed)
- **Server Components**: Leverage Next.js App Router for data fetching where possible

### Testing Strategy (Speed-First)
- **Unit Tests**: Critical business logic only (services, calculations)
- **Integration Tests**: Key API endpoints + database operations
- **E2E Tests**: Defer until core features are stable
- **Coverage Target**: 60% initially, focus on high-value paths

### Deployment Strategy (Speed-First)
- **Phase 1 (Now)**: Local development with `docker-compose` (Postgres, Redis, services)
- **Phase 2 (Later)**: AWS ECS/Fargate when ready to scale
- **Infrastructure as Code**: Prepare CDK structure, implement when needed

---

## Implementation Order (Speed-First Approach)

### 🎯 Objective 0: Supabase Setup (30 minutes)

**0.1 Create Supabase Project**
1. Go to [supabase.com](https://supabase.com)
2. Sign up / Log in
3. Create new project:
   - Project name: `periscope` or `periscope-dev`
   - Database password: Generate a strong password (save it!)
   - Region: Choose closest to you (e.g., `us-east-1`)
   - Pricing: Free tier is sufficient for development

**0.2 Get API Credentials**
Once project is created:
1. Go to **Settings** → **API**
2. Copy the following:
   - **Project URL** → `SUPABASE_URL`
   - **anon public** key → `SUPABASE_ANON_KEY`
   - **service_role** key → `SUPABASE_SERVICE_KEY` (keep secret!)

**0.3 Configure Database Extensions**
1. Go to **Database** → **Extensions**
2. Enable the following extensions:
   - ✅ `pg_stat_statements` - Query performance monitoring
   - ✅ `timescaledb` - Time-series data (for historical options data)
   - ✅ `pg_cron` - Scheduled jobs (optional, for later)

**0.4 Create Initial Schema (Optional - can defer to Week 2)**
Go to **SQL Editor** and create tables:
```sql
-- Options contracts table
CREATE TABLE options_contracts (
  id BIGSERIAL PRIMARY KEY,
  ticker TEXT NOT NULL,
  contract_type TEXT NOT NULL, -- 'call' or 'put'
  strike_price NUMERIC(10, 2) NOT NULL,
  expiration_date DATE NOT NULL,
  underlying_ticker TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Options quotes (time-series data)
CREATE TABLE options_quotes (
  id BIGSERIAL PRIMARY KEY,
  ticker TEXT NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL,
  bid NUMERIC(10, 4),
  ask NUMERIC(10, 4),
  implied_volatility NUMERIC(10, 6),
  delta NUMERIC(10, 6),
  gamma NUMERIC(10, 6),
  theta NUMERIC(10, 6),
  vega NUMERIC(10, 6),
  open_interest BIGINT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Convert to TimescaleDB hypertable for better performance
SELECT create_hypertable('options_quotes', 'timestamp', if_not_exists => TRUE);

-- Create indexes
CREATE INDEX idx_options_contracts_ticker ON options_contracts(ticker);
CREATE INDEX idx_options_contracts_underlying ON options_contracts(underlying_ticker);
CREATE INDEX idx_options_contracts_expiration ON options_contracts(expiration_date);
CREATE INDEX idx_options_quotes_ticker ON options_quotes(ticker);
CREATE INDEX idx_options_quotes_timestamp ON options_quotes(timestamp DESC);
```

**0.5 Update Environment Variables**
Create `.env` file in project root:
```bash
# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key

# Massive API
MASSIVE_API_KEY=your-massive-api-key

# Backend
PORT=8080
GIN_MODE=debug
```

**Critical Info:**
- ✅ Supabase free tier: 500MB database, 2GB bandwidth/month, 50K monthly active users
- ✅ PostgreSQL 15 with TimescaleDB support
- ✅ Built-in authentication (use later for user login)
- ✅ Real-time subscriptions (use for live options updates)

---

### 🎯 Objective 1: Go Backend Core Structure (Days 1-2)

**1.1 Project Initialization**
```bash
mkdir backend-go && cd backend-go
go mod init github.com/yourusername/periscope-backend
```

**1.2 Create Directory Structure**
```
backend-go/
├── cmd/api/main.go
├── internal/
│   ├── api/
│   │   ├── handlers/options.go
│   │   ├── middleware/cors.go
│   │   ├── middleware/logger.go
│   │   └── router.go
│   ├── services/options_service.go
│   └── models/
│       ├── options.go
│       ├── greeks.go
│       └── response.go
├── pkg/
│   ├── massive/client.go
│   └── errors/errors.go
├── config/config.go
├── .env.example
├── go.mod
├── Makefile
└── README.md
```

**1.3 Install Core Dependencies**
```bash
go get github.com/gin-gonic/gin
go get github.com/spf13/viper
go get github.com/stretchr/testify
go get golang.org/x/time/rate           # rate limiting
go get github.com/jackc/pgx/v5          # PostgreSQL driver for SQLC
go get github.com/jackc/pgx/v5/pgxpool  # Connection pooling
go get github.com/supabase-community/supabase-go  # Supabase Go client (optional)
```

**1.4 Configuration Setup**
- Create `config/config.go` with Viper
- Load from `.env` file
- Environment variables:
  - `MASSIVE_API_KEY` - Massive.com API key
  - `SUPABASE_URL` - Supabase project URL
  - `SUPABASE_ANON_KEY` - Supabase anonymous key
  - `SUPABASE_SERVICE_KEY` - Supabase service role key (for backend)
  - `PORT` - Server port (default: 8080)
  - `GIN_MODE` - Gin mode (debug/release)

**1.5 Setup Supabase Database Connection**
```go
// pkg/database/supabase.go
package database

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
    Pool *pgxpool.Pool
}

func NewSupabaseDB(connectionString string) (*DB, error) {
    config, err := pgxpool.ParseConfig(connectionString)
    if err != nil {
        return nil, err
    }

    // Connection pool settings
    config.MaxConns = 10
    config.MinConns = 2

    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, err
    }

    // Test connection
    if err := pool.Ping(context.Background()); err != nil {
        return nil, err
    }

    return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
    db.Pool.Close()
}
```

**Connection String Format:**
```
postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
```

Get from Supabase: **Settings** → **Database** → **Connection string** → **Golang**

**1.6 Port Massive API Client**
- Create `pkg/massive/client.go`
- Port structs from Rust (OptionContract, Greeks, etc.)
- Add rate limiting (10 req/sec)
- Add basic retry logic (3 attempts with exponential backoff)

**1.7 Create Basic REST API**
- Initialize Gin router in `main.go`
- Single endpoint: `GET /api/v1/options/:ticker`
- CORS middleware (allow `http://localhost:3000` for Next.js)
- Structured logging middleware
- Database health check endpoint: `GET /api/v1/health`

**1.8 Testing**
- Unit tests for Massive client
- Test with real API key (mark with build tag `// +build integration`)
- Test database connection
- `make test` target in Makefile

**Critical Files:**
- `/Users/aaronbengo/Documents/github/optionsPricing/backend-go/cmd/api/main.go`
- `/Users/aaronbengo/Documents/github/optionsPricing/backend-go/pkg/massive/client.go`
- `/Users/aaronbengo/Documents/github/optionsPricing/backend-go/internal/api/router.go`

---

### 🎯 Objective 2: Next.js Frontend Setup (Days 2-3)

**2.1 Initialize Next.js Project**
```bash
npx create-next-app@latest frontend --typescript --tailwind --app --no-src-dir
cd frontend
```

**2.2 Install Dependencies**
```bash
npm install @tanstack/react-query axios
npm install @supabase/supabase-js  # Supabase client
npm install recharts  # 2D charts
npm install plotly.js plotly.js-dist-min react-plotly.js  # 3D surfaces
npm install -D @types/react-plotly.js
npm install zustand  # client state (if needed)
```

**2.3 Project Structure**
```
frontend/
├── app/
│   ├── layout.tsx              # Root layout with providers
│   ├── page.tsx                # Dashboard home
│   ├── providers.tsx           # TanStack Query provider
│   └── dashboard/
│       ├── page.tsx            # Main dashboard (Server Component)
│       └── components/
│           ├── OptionsTable.tsx      # Client Component
│           ├── GreeksDisplay.tsx     # Client Component
│           └── IVChart.tsx           # Client Component
├── lib/
│   ├── api/
│   │   ├── client.ts           # Axios instance
│   │   └── queries.ts          # TanStack Query hooks
│   └── utils.ts
├── components/
│   └── ui/                     # shadcn/ui components (later)
├── .env.local
└── next.config.mjs
```

**2.4 Setup API Clients**

**Go Backend API Client:**
```typescript
// lib/api/client.ts
import axios from 'axios';

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
});

// lib/api/queries.ts
import { useQuery } from '@tanstack/react-query';

export function useOptionsChain(ticker: string) {
  return useQuery({
    queryKey: ['options', ticker],
    queryFn: async () => {
      const { data } = await apiClient.get(`/options/${ticker}`);
      return data;
    },
    refetchInterval: 30000, // Refetch every 30s
  });
}
```

**Supabase Client (for direct database access):**
```typescript
// lib/supabase/client.ts
import { createClient } from '@supabase/supabase-js';

const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL!;
const supabaseAnonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!;

export const supabase = createClient(supabaseUrl, supabaseAnonKey);

// Example: Query options quotes directly from Supabase
export async function getHistoricalQuotes(ticker: string, startDate: Date) {
  const { data, error } = await supabase
    .from('options_quotes')
    .select('*')
    .eq('ticker', ticker)
    .gte('timestamp', startDate.toISOString())
    .order('timestamp', { ascending: false });

  if (error) throw error;
  return data;
}
```

**Note:** For Week 1, we'll primarily use the Go backend API. Direct Supabase queries can be added later for real-time features.

**2.5 Create Dashboard Components**
- `app/layout.tsx`: Setup TanStack Query provider
- `app/dashboard/page.tsx`: Main dashboard with ticker selector
- `app/dashboard/components/OptionsTable.tsx`: Display options chain with strikes, IVs, Greeks
- `app/dashboard/components/IVChart.tsx`: Recharts line chart for IV vs strike

**2.6 Environment Variables**
```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
```

**2.7 Testing**
- Run `npm run dev`
- Navigate to `http://localhost:3000/dashboard`
- Select ticker (AAPL)
- Verify data flows from Go backend

**Critical Files:**
- `/Users/aaronbengo/Documents/github/optionsPricing/frontend/app/layout.tsx`
- `/Users/aaronbengo/Documents/github/optionsPricing/frontend/lib/api/client.ts`
- `/Users/aaronbengo/Documents/github/optionsPricing/frontend/app/dashboard/components/OptionsTable.tsx`

---

### 🎯 Objective 3: Massive API Integration (Day 1 - Done in Objective 1)

**Covered in Objective 1.5 - No separate phase needed**

Enhancements already included:
- ✅ Rate limiting (token bucket, 10 req/sec)
- ✅ Retry logic (3 attempts, exponential backoff)
- ✅ Error handling with custom error types
- ✅ Structured logging

**Future Enhancements (Deferred):**
- Circuit breaker (when scaling)
- Pagination handling (follow `next_url`)
- Request coalescing (dedupe concurrent requests)

---

### 🎯 Objective 4: Request/Response Objects (Day 2 - Done in Objective 1)

**Already covered in Objective 1.5 - Port data models from Rust**

**Data Structures:**
```go
// internal/models/options.go
type OptionsChainResponse struct {
    Status    string           `json:"status"`
    RequestID string           `json:"request_id"`
    Results   []OptionContract `json:"results"`
    NextURL   *string          `json:"next_url,omitempty"`
}

type OptionContract struct {
    Details         ContractDetails `json:"details"`
    Greeks          Greeks          `json:"greeks"`
    ImpliedVol      *float64        `json:"implied_volatility"`
    OpenInterest    *int64          `json:"open_interest"`
    LastQuote       LastQuote       `json:"last_quote"`
    LastTrade       LastTrade       `json:"last_trade"`
    Day             DayBar          `json:"day"`
    UnderlyingAsset UnderlyingAsset `json:"underlying_asset"`
}

// ... (mirror all Rust structs)
```

**Validation:**
- Use `binding` tags for Gin
```go
type GetOptionsRequest struct {
    Ticker string `uri:"ticker" binding:"required,alphanum,min=1,max=10"`
    Limit  int    `form:"limit" binding:"omitempty,min=1,max=1000"`
}
```

**API Versioning:** URL path `/api/v1/*`

**OpenAPI Docs (Optional for now - defer to later):**
- Add swag annotations later
- Generate with `swag init`

---

### 🎯 Objective 5: Basic CI/CD (Days 4-5)

**5.1 GitHub Actions - CI Workflow**

Create `.github/workflows/ci.yml`:
```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Install dependencies
        working-directory: backend-go
        run: go mod download
      - name: Run tests
        working-directory: backend-go
        run: go test -v ./...
      - name: Run linter
        working-directory: backend-go
        run: |
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
          golangci-lint run

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - name: Install dependencies
        working-directory: frontend
        run: npm ci
      - name: Run lint
        working-directory: frontend
        run: npm run lint
      - name: Build
        working-directory: frontend
        run: npm run build
```

**5.2 Docker Setup (for local development)**

Create `docker-compose.yml` at project root:
```yaml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  backend:
    build:
      context: ./backend-go
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - SUPABASE_URL=${SUPABASE_URL}
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
      - SUPABASE_SERVICE_KEY=${SUPABASE_SERVICE_KEY}
      - REDIS_URL=redis://redis:6379
      - MASSIVE_API_KEY=${MASSIVE_API_KEY}
      - PORT=8080
      - GIN_MODE=debug
    depends_on:
      - redis
    env_file:
      - .env

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
      - NEXT_PUBLIC_SUPABASE_URL=${SUPABASE_URL}
      - NEXT_PUBLIC_SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
    depends_on:
      - backend
    env_file:
      - .env

volumes:
  redis_data:
```

**Note:** PostgreSQL is managed by Supabase, so no local database container needed!

**5.3 Makefiles for Automation**

**backend-go/Makefile:**
```makefile
.PHONY: build run test lint

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

docker-build:
	docker build -t periscope-backend .
```

**frontend/Makefile:**
```makefile
.PHONY: install dev build lint

install:
	npm install

dev:
	npm run dev

build:
	npm run build

lint:
	npm run lint
```

**5.4 Testing Strategy**
- **Backend**: Unit tests for `pkg/massive`, `internal/services`
- **Frontend**: Component tests deferred (focus on functionality first)
- **Integration**: Manual testing with real API
- **Coverage**: Measure but don't enforce threshold yet

**Critical Files:**
- `/Users/aaronbengo/Documents/github/optionsPricing/.github/workflows/ci.yml`
- `/Users/aaronbengo/Documents/github/optionsPricing/docker-compose.yml`

---

## Verification & Testing Plan

### Manual End-to-End Testing

**Step 1: Test Go Backend Standalone**
```bash
cd backend-go

# Setup environment
cp .env.example .env
# Edit .env with MASSIVE_API_KEY

# Install dependencies
go mod download

# Run tests
make test

# Start server
make run

# In another terminal, test API
curl http://localhost:8080/api/v1/options/AAPL | jq
```

**Expected Output:**
```json
{
  "status": "OK",
  "request_id": "...",
  "results": [
    {
      "details": {
        "ticker": "O:AAPL240119C00150000",
        "contract_type": "call",
        "strike_price": 150.0,
        ...
      },
      "greeks": {
        "delta": 0.52,
        "gamma": 0.03,
        ...
      },
      ...
    }
  ]
}
```

---

**Step 2: Test Next.js Frontend**
```bash
cd frontend

# Setup environment
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" > .env.local

# Install dependencies
npm install

# Start dev server
npm run dev
```

**Manual Testing:**
1. Open `http://localhost:3000/dashboard` in browser
2. Select ticker "AAPL" from dropdown
3. Verify options chain table displays:
   - Strike prices
   - Implied volatility
   - Greeks (Delta, Gamma, Theta, Vega)
   - Bid/Ask spreads
4. Check browser console for errors
5. Verify data refreshes every 30 seconds
6. Test with different tickers: TSLA, SPY, MSFT

---

**Step 3: Test Full Stack with Docker Compose**
```bash
# From project root
docker-compose up --build

# Wait for services to start (~30 seconds)
# Check logs for "Server started on :8080"

# Test backend
curl http://localhost:8080/api/v1/options/SPY | jq

# Open frontend
open http://localhost:3000/dashboard
```

---

**Step 4: Test CI/CD Pipeline**
```bash
# Create feature branch
git checkout -b feature/test-pipeline

# Make a small change (add comment)
echo "// Test change" >> backend-go/cmd/api/main.go

# Commit and push
git add .
git commit -m "test: verify CI pipeline"
git push origin feature/test-pipeline

# Create PR on GitHub
# Verify CI workflow runs and passes
```

**Expected CI Results:**
- ✅ Backend tests pass
- ✅ Backend linter passes
- ✅ Frontend build succeeds
- ✅ Frontend linter passes

---

## Success Criteria

- [ ] Go API responds to `GET /api/v1/options/AAPL`
- [ ] Next.js dashboard displays options chain with Greeks
- [ ] CI pipeline passes on every commit
- [ ] `docker-compose up` runs full stack locally

---

## Next Steps After Week 1

Once the foundational 5 objectives are complete, prioritize:

1. **Database Persistence (Week 2)**
   - Supabase is already set up! ✅
   - Define detailed schema for options, quotes, trades (expand from basic schema)
   - Implement SQLC queries for type-safe database access
   - Repository layer for data access
   - Add migrations using golang-migrate or Supabase migrations

2. **Caching Layer (Week 2)**
   - Integrate Redis
   - Cache Massive API responses (5-minute TTL)
   - Cache-aside pattern implementation

3. **Enhanced Error Handling (Week 2)**
   - Circuit breaker for Massive API
   - Structured error responses
   - Sentry/error tracking integration

4. **Advanced Analytics (Week 3+)**
   - Port Rust analytics engine (Greeks calculations)
   - IV surface construction
   - Skew/smile analysis

5. **AWS Deployment (Week 3+)**
   - CDK infrastructure setup
   - ECS task definitions
   - RDS database setup
   - CI/CD deployment pipeline
