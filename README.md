# Periscope

A high-performance Rust platform for options volatility surface analysis, relative value identification, and systematic trading.

## Overview

Periscope is a quantitative options analytics engine designed to:

1. **Extract volatility skew and smile** from live options chains across all strikes and expirations
2. **Construct arbitrage-free volatility surfaces** using SVI/SSVI parameterization
3. **Calibrate Heston stochastic volatility models** to observed market data
4. **Identify relative value opportunities** by comparing model-implied vs market-implied volatility
5. **Visualize mispricing** as a temperature-mapped 3D surface (blue = cheap, red = rich)
6. **Generate systematic trade signals** for skew trades, ratio spreads, and butterflies

The backend is built in Rust for sub-millisecond latency, enabling real-time analysis of rapidly changing options markets.

## Quantitative Strategy

Periscope implements a volatility relative value framework:

- **Skew Analysis**: Compare delta-equivalent puts vs calls to identify directional mispricing
- **Smile Analysis**: Analyze strike-space curvature to identify tail risk mispricing
- **Surface Construction**: Build complete 3D IV surfaces across all expirations
- **Heston Calibration**: Fit stochastic volatility model to extract fair-value IVs
- **Relative Value**: Trade options that deviate significantly from model fair value

See [`goal.md`](goal.md) for comprehensive quantitative methodology and [`plan.md`](plan.md) for project roadmap.

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Market Data    │────▶│    Periscope    │────▶│    Frontend     │
│  (Massive API)  │     │  (Rust Engine)  │     │   (WebSocket)   │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                               │
                               ▼
                        ┌─────────────────┐
                        │  Analytics      │
                        │  ├─ Skew        │
                        │  ├─ Smile       │
                        │  ├─ Surface     │
                        │  ├─ Heston      │
                        │  └─ RelValue    │
                        └─────────────────┘
```

**Why Rust?**
- Zero garbage collection pauses for predictable latency
- Memory safety without runtime overhead
- Async I/O for handling thousands of concurrent connections
- Performance comparable to C++ with modern tooling

## Project Structure

```
periscope/
├── src/
│   ├── lib.rs              # Library root with public exports
│   ├── error.rs            # Centralized error types (thiserror)
│   ├── bin/                # Binary entry points
│   │   └── greeks_test.rs  # CLI for options chain fetching
│   ├── client/             # External API clients
│   │   └── massive.rs      # Massive API integration
│   ├── config/             # Environment configuration
│   ├── models/             # Data structures
│   │   ├── greeks.rs       # Greeks (Δ, Γ, Θ, V)
│   │   └── options.rs      # Option contracts, quotes, trades
│   └── services/           # Business logic (future: pricing, analytics)
├── examples/               # Usage examples for library consumers
├── tests/                  # Integration tests
├── py_quick_test/          # Python prototypes (gitignored)
├── goal.md                 # Quantitative methodology & technical specs
├── plan.md                 # Project roadmap & progress tracking
└── Makefile                # Build automation
```

## Documentation

| Document | Description |
|----------|-------------|
| [`goal.md`](goal.md) | Comprehensive quantitative methodology: skew/smile extraction, SVI parameterization, Heston calibration, relative value framework, trade signals |
| [`plan.md`](plan.md) | Project roadmap with phased implementation plan and progress tracking |
| [`idea.md`](idea.md) | Original strategy concept and trading rationale |

## Prerequisites

- [Rust](https://rustup.rs/) (1.70+)
- Make
- Massive API key ([massive.com](https://massive.com))

## Getting Started

### 1. Clone and configure

```bash
git clone git@github.com:aaronbengochea/Periscope.git
cd Periscope

# Copy environment template and add your API key
cp .env.example .env
# Edit .env with your MASSIVE_API_KEY
```

### 2. Build

```bash
# Debug build (fast compile, for development)
make build

# Release build (optimized, for production)
make build-release
```

### 3. Run

```bash
# Run with default settings (AAPL, 10 contracts)
make run-release

# Run with custom arguments
make run-release ARGS='--ticker TSLA --limit 20'

# Run example
make example
```

### 4. Test & Lint

```bash
make test    # Run all tests
make fmt     # Format code
make lint    # Run clippy linter
```

### 5. Cleanup

```bash
make clean   # Remove build artifacts
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make build` | Build debug binary |
| `make build-release` | Build optimized release binary |
| `make run` | Run debug binary |
| `make run-release` | Run release binary |
| `make example` | Run basic_usage example |
| `make test` | Run all tests |
| `make fmt` | Format code |
| `make lint` | Run clippy linter |
| `make clean` | Remove build artifacts |

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `MASSIVE_API_KEY` | Massive API key for options data | Yes |
| `MASSIVE_BASE_URL` | API base URL | No (defaults to `https://api.massive.com/v3`) |
| `RUST_LOG` | Logging level (e.g., `info`, `debug`) | No |

## Roadmap

See [`plan.md`](plan.md) for detailed roadmap. High-level phases:

1. **Data Infrastructure** — Options chain ingestion, storage, normalization
2. **Skew Engine** — Delta-space analysis, risk reversal calculation
3. **Smile Engine** — Strike-space analysis, SVI fitting
4. **Surface Builder** — SSVI parameterization, arbitrage-free interpolation
5. **Heston Calibration** — Stochastic vol model fitting
6. **Relative Value** — Model vs market comparison, mispricing detection
7. **Visualization** — 3D temperature-mapped surface (WebGL)
8. **Trade Signals** — Systematic signal generation
9. **Execution** — Broker integration, automated trading

## License

MIT
