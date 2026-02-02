# Periscope

A high-performance Rust backend for real-time options analysis, visualization, and recommendations.

## Overview

Periscope is designed to power millisecond-interval options data processing, enabling:

- **Real-time options chain analysis** with Greeks (Delta, Gamma, Theta, Vega)
- **Low-latency data streaming** to frontend clients via WebSocket
- **Options strategy recommendations** based on market conditions
- **Visual analytics** through a modern, responsive frontend

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Market Data    │────▶│    Periscope    │────▶│    Frontend     │
│  (Massive API)  │     │  (Rust Backend) │     │   (WebSocket)   │
└─────────────────┘     └─────────────────┘     └─────────────────┘
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
│   ├── lib.rs              # Library root
│   ├── error.rs            # Centralized error types
│   ├── bin/                # Binary entry points
│   ├── client/             # API clients (Massive, etc.)
│   ├── config/             # Configuration management
│   ├── models/             # Data structures
│   └── services/           # Business logic
├── examples/               # Usage examples
├── tests/                  # Integration tests
├── py_quick_test/          # Python prototypes
└── Makefile                # Build commands
```

## Prerequisites

- [Rust](https://rustup.rs/) (1.70+)
- Make

## Getting Started

### 1. Clone and configure

```bash
git clone git@github.com:aaronbengochea/Periscope.git
cd Periscope

# Copy environment template and add your API key
cp .env.example .env
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
| `MASSIVE_API_KEY` | Massive API key | Yes |
| `MASSIVE_BASE_URL` | API base URL | No (defaults to `https://api.massive.com/v3`) |

## License

MIT
