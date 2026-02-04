# Periscope

The goal of this project is to build a high performance/reliability Rust based platform for agentic enabled options volatility surface analysis, vol relative value identification, and deterministic programatic systematic trading. 

If you feel there is something wrong with the world, instead of complaining, you should aim to fix it. Capital is but a tool and it earns one of the, if not the, highest PCA weights in our current existing world relative to solving real problems. This project acts as my manifesto to the world, fight the good fight, earn your stake in this world, be the agent you were meant to be.

A primal instinct: Risk and Reward. All so that we can free ourselves from the shackels of the modern world and enable survivability. Proof that past, present, and future are initmately linked. Freedom is stackable, and time is not to be squandered or wasted, especially when you were gifted with the ability and belief that you could one day change the world for the better. 

_**SOLUS ET INTREPIDUS**_

_**ATOP FURTIM VIGILANS**_

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

See [`goal.md`](goal.md) for comprehensive quantitative methodology, [`quantitative_plan.md`](quantitative_plan.md) for analytics roadmap, and [`infrastructure_plan.md`](infrastructure_plan.md) for systems architecture.

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

## Current Project Structure

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
├── quantitative_plan.md    # Analytics implementation roadmap
├── infrastructure_plan.md  # Systems architecture & infrastructure roadmap
└── Makefile                # Build automation
```

## Documentation

| Document | Description |
|----------|-------------|
| [`goal.md`](goal.md) | Comprehensive quantitative methodology: skew/smile extraction, SVI parameterization, Heston calibration, relative value framework, trade signals |
| [`quantitative_plan.md`](quantitative_plan.md) | Analytics implementation roadmap: IV extraction, skew/smile engines, surface construction, Heston calibration, relative value, trade signals |
| [`infrastructure_plan.md`](infrastructure_plan.md) | Systems architecture: Rust services, AWS CDK infrastructure, database schemas, API design, frontend components |
| [`original_idea.md`](original_idea.md) | Original strategy concept and trading rationale |

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

See [`infrastructure_plan.md`](infrastructure_plan.md) for systems architecture and [`quantitative_plan.md`](quantitative_plan.md) for analytics implementation. High-level phases:

1. **Data Infrastructure** — Options chain ingestion, storage, normalization
2. **Skew Engine** — Delta-space analysis, risk reversal calculation
3. **Smile Engine** — Strike-space analysis, SVI fitting
4. **Surface Builder** — SSVI parameterization, arbitrage-free interpolation
5. **Heston Calibration** — Stochastic vol model fitting
6. **Relative Value** — Model vs market comparison, mispricing detection
7. **Visualization** — 3D temperature-mapped surface (WebGL)
8. **Trade Signals** — Systematic signal generation
9. **Execution** — Broker integration, automated trading

