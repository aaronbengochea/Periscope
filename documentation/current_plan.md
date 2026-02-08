# Periscope - Current Implementation Status & Roadmap

## 🎯 Project Overview

Periscope is transitioning to a **dual-backend architecture**:
- **Go Backend**: General API operations, data aggregation, user management
- **Next.js Frontend**: Modern React-based UI with real-time data
- **Rust Engine** (Future): High-performance analytics and order execution

---

## ✅ Completed (As of Feb 7, 2026)

### Backend Infrastructure
- ✅ **Go Backend with Gin framework**
  - RESTful API with CORS middleware
  - Structured logging
  - Health check endpoint
  - Rate limiting (10 req/sec)

- ✅ **Massive API Integration**
  - Full options chain endpoint (`/snapshot/options/{ticker}`)
  - Automatic pagination (handles 20+ pages, 5000+ contracts)
  - Stock price fetching via unified snapshot
  - Price injection into contracts

- ✅ **Dual-Endpoint Strategy for Enriched Data**
  - Initial load: Fast options chain snapshot (all contracts)
  - Expiration selection: Detailed unified snapshot (50-100 contracts)
  - Smart batching (up to 250 contracts per request)
  - Data merging: Enriched data overlays basic data
  - **Greeks now available**: Delta, Gamma, Theta, Vega

- ✅ **Docker Infrastructure**
  - docker-compose setup for local development
  - Backend and frontend containers
  - Makefile commands for easy management

### Frontend Infrastructure
- ✅ **Next.js 15 with TypeScript**
  - App Router architecture
  - TanStack Query for data fetching
  - Auto-refresh every 30 seconds

- ✅ **Options Chain Display**
  - 13-column layout (Last, Bid, Bid Sz, Ask, Ask Sz, IV, Greeks, OI, Chg%, Vol)
  - Current price indicator inline with search (top right)
  - Expiration dropdown with formatted dates
  - Days to expiry calculation
  - Strike highlighting (ATM detection)
  - CALLS/PUTS header row

- ✅ **Smart Data Handling**
  - Filters contracts by selected expiration
  - Only shows future/current expirations
  - Graceful handling of missing data (shows "-")
  - Subscription limitation notice

### Data Models
- ✅ **Complete type definitions** (Go + TypeScript)
  - OptionContract with all fields
  - Greeks (Delta, Gamma, Theta, Vega, Rho)
  - LastQuote (Bid, Ask, sizes)
  - DayBar (OHLCV data)
  - Session (Change $, Change %, volume)
  - ContractDetails
  - UnderlyingAsset

---

## 📊 Current Data Availability

### Available Now (Free Tier)
- ✅ **Greeks**: Delta, Gamma, Theta, Vega (via unified snapshot)
- ✅ **Implied Volatility**
- ✅ **Open Interest**
- ✅ **Day OHLCV** (Open, High, Low, Close, Volume)
- ✅ **Contract Details** (Strike, Expiration, Type)
- ✅ **Stock Price** (underlying asset)

### Requires Subscription Upgrade
- ❌ **Bid/Ask Quotes**: Requires quotes subscription
- ❌ **Bid/Ask Sizes**: Requires quotes subscription
- ❌ **Session Data**: Change $, Change % (empty in response)
- ❌ **Last Trade**: Price, Size (not populated)

---

## 🚧 In Progress / Next Steps

### Immediate (Week 1-2)
1. **UI Polish**
   - [ ] Add loading states for detail fetch
   - [ ] Error handling for failed detail fetches
   - [ ] Optimize layout for different screen sizes
   - [ ] Add tooltips for Greeks explanations

2. **Performance Optimization**
   - [ ] Implement request debouncing
   - [ ] Add caching for enriched data (30s TTL)
   - [ ] Reduce unnecessary re-renders
   - [ ] Profile API call patterns

3. **Data Quality**
   - [ ] Handle edge cases (missing contracts, API errors)
   - [ ] Validate data before display
   - [ ] Add data freshness indicators

### Short Term (Weeks 2-4)

4. **Supabase Database Integration**
   - [ ] Set up Supabase project
   - [ ] Create schema for options contracts
   - [ ] Create schema for historical quotes (TimescaleDB)
   - [ ] Implement SQLC for type-safe queries
   - [ ] Add repository layer
   - [ ] Database connection pooling

5. **Caching Layer**
   - [ ] Integrate Redis
   - [ ] Cache Massive API responses (5-min TTL)
   - [ ] Cache enriched contract data
   - [ ] Implement cache-aside pattern
   - [ ] Add cache invalidation logic

6. **Testing & CI/CD**
   - [ ] Unit tests for backend services
   - [ ] Integration tests for API endpoints
   - [ ] GitHub Actions CI workflow
   - [ ] Automated testing on PRs
   - [ ] Code coverage reporting

### Medium Term (Weeks 4-8)

7. **Advanced Features**
   - [ ] Historical data storage
   - [ ] IV surface visualization (3D charts)
   - [ ] Skew/smile analysis
   - [ ] Volatility term structure
   - [ ] Multiple ticker comparison

8. **User Features**
   - [ ] Authentication (Supabase Auth)
   - [ ] User portfolios
   - [ ] Watchlists
   - [ ] Custom alerts
   - [ ] Trade journal

9. **Analytics Engine**
   - [ ] Black-Scholes pricing
   - [ ] Greeks calculation (independent verification)
   - [ ] IV calculation from prices
   - [ ] Risk metrics (VaR, Greeks exposure)

### Long Term (Weeks 8+)

10. **Production Deployment**
    - [ ] AWS ECS/Fargate setup
    - [ ] CDK infrastructure as code
    - [ ] Production database (Supabase production tier)
    - [ ] Monitoring & alerting (Sentry, DataDog)
    - [ ] Load balancing
    - [ ] Auto-scaling

11. **Rust Integration**
    - [ ] High-performance Greeks engine
    - [ ] Real-time pricing
    - [ ] Order execution simulation
    - [ ] Backtesting engine

---

## 🏗️ Architecture Decisions

### Backend Stack
- **Framework**: Gin (Go web framework)
- **Database**: Supabase (PostgreSQL with TimescaleDB) - *Not yet integrated*
- **ORM**: SQLC (type-safe SQL code generation) - *Not yet integrated*
- **Cache**: Redis (planned)
- **API**: RESTful with JSON

### Frontend Stack
- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript
- **State**: TanStack Query (server state) + Zustand (client state, minimal)
- **Styling**: Tailwind CSS
- **Charts**: Recharts (2D) + Plotly.js (3D, planned)

### Data Strategy
- **Primary Source**: Massive API (options chain + unified snapshot)
- **Storage**: Supabase (future - for historical data)
- **Caching**: Redis (future - for real-time performance)
- **Real-time**: Supabase subscriptions (future)

---

## 📁 Project Structure

```
optionsPricing/
├── backend-go/              # Go backend API
│   ├── cmd/api/            # Application entry point
│   ├── internal/
│   │   ├── api/            # HTTP handlers, routes, middleware
│   │   └── models/         # Data structures
│   └── pkg/
│       ├── massive/        # Massive API client
│       ├── database/       # Database connections (future)
│       └── errors/         # Error handling
│
├── frontend/               # Next.js frontend
│   ├── app/               # Next.js App Router
│   ├── components/        # React components
│   └── lib/               # API clients, utilities
│
├── docker-compose.yml     # Local development setup
├── Makefile              # Project automation
└── current_plan.md       # This file
```

---

## 🔑 Key Files

### Backend
- `backend-go/cmd/api/main.go` - Application entry point
- `backend-go/pkg/massive/client.go` - Massive API integration
  - `GetOptionsChain()` - Fetch options with pagination
  - `GetContractDetails()` - Fetch enriched data via unified snapshot
  - `GetStockPrice()` - Fetch underlying stock price
- `backend-go/internal/api/router.go` - Route definitions
- `backend-go/internal/api/handlers/options.go` - Options endpoints
  - `GET /api/v1/options/:ticker` - Get options chain
  - `POST /api/v1/options/details` - Get enriched contract details
- `backend-go/internal/models/options.go` - Data models

### Frontend
- `frontend/app/page.tsx` - Main dashboard with search & expiration selection
- `frontend/components/OptionsChain.tsx` - Options chain table (13 columns)
- `frontend/components/ExpirationDropdown.tsx` - Expiration selector
- `frontend/lib/api.ts` - API client functions
- `frontend/lib/dateUtils.ts` - Date formatting utilities

---

## 🎓 Lessons Learned

1. **Dual-Endpoint Strategy**: Combining fast initial load with on-demand enrichment provides best UX
2. **URL Construction**: Be careful with baseURL - avoid double /v3 paths
3. **Subscription Tiers**: Free tier provides Greeks but not quotes - plan accordingly
4. **Date Parsing**: Use local time component parsing to avoid timezone bugs
5. **Data Merging**: Merge enriched data carefully to avoid overwriting good data with undefined
6. **Pagination**: Massive API requires manual next_url following and API key injection

---

## 🚀 How to Run

### Prerequisites
- Go 1.22+
- Node.js 20+
- Docker Desktop
- Massive API key

### Quick Start
```bash
# 1. Set environment variables
cp .env.example .env
# Edit .env with your MASSIVE_API_KEY

# 2. Start all services
make docker-up

# 3. Access the app
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# Health check: http://localhost:8080/health

# 4. Stop services
make docker-down
```

### Development Mode
```bash
# Backend only (without Docker)
cd backend-go
make run

# Frontend only (without Docker)
cd frontend
npm run dev
```

---

## 📝 Notes

- **API Rate Limits**: Massive API limited to 10 req/sec (enforced client-side)
- **Pagination**: Auto-fetches up to 20 pages (5,000 contracts) - safety limit
- **Greeks Source**: Unified snapshot endpoint provides fully populated Greeks
- **Quotes**: Requires subscription upgrade - UI gracefully handles missing data
- **Time Zone**: Application assumes America-only for date calculations

---

## 🐛 Known Issues

1. **Missing Quote Data**: Bid/Ask not available without subscription upgrade (expected)
2. **Empty Session**: Session data structure exists but fields are null (subscription issue)
3. **Ticker Missing**: Contract details from unified snapshot don't include ticker field

---

## 💡 Future Enhancements

- Real-time updates via WebSocket or Supabase subscriptions
- Historical IV tracking and charting
- Options strategy builder (spreads, butterflies, condors)
- Portfolio Greek calculations
- What-if scenario analysis
- Mobile-responsive design
- Dark/light theme toggle
- Export to CSV/Excel
- Integration with brokers for live trading

---

**Last Updated**: February 7, 2026
**Status**: Active Development
**Version**: 0.2.0 (Dual-Endpoint Strategy)
