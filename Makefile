# ==============================================================================
# Periscope Makefile
# ==============================================================================
# This Makefile provides convenient shortcuts for common Rust development tasks.
# Run `make` or `make help` to see all available commands.
# ==============================================================================

.PHONY: build build-release run run-release check test fmt lint clean help run_go_front docker-up docker-down docker-logs docker-build docker-restart

# ==============================================================================
# DEFAULT TARGET
# ==============================================================================
# Running `make` with no arguments displays this help menu.
# Lists all available commands with brief descriptions.
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build:"
	@echo "  build          Build debug binary"
	@echo "  build-release  Build optimized release binary"
	@echo "  check          Check for errors without building"
	@echo ""
	@echo "Run:"
	@echo "  run            Run debug binary"
	@echo "  run-release    Run release binary"
	@echo "  example        Run basic_usage example"
	@echo "  run_go_front   Run Go backend + Next.js frontend (topological order)"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   Build Docker images"
	@echo "  docker-up      Start all services with Docker Compose"
	@echo "  docker-down    Stop all services"
	@echo "  docker-logs    View logs from all services"
	@echo "  docker-restart Restart all services"
	@echo ""
	@echo "Test & Lint:"
	@echo "  test           Run all tests"
	@echo "  fmt            Format code"
	@echo "  lint           Run clippy linter"
	@echo ""
	@echo "Other:"
	@echo "  clean          Remove build artifacts"
	@echo ""
	@echo "Examples:"
	@echo "  make run-release ARGS='--ticker TSLA --limit 20'"

# ==============================================================================
# BUILD TARGETS
# ==============================================================================

# build: Compiles the project in debug mode.
#   - Fast compilation, slower runtime performance
#   - Includes debug symbols for debugging with gdb/lldb
#   - No optimizations applied
#   - Output binary: target/debug/greeks_test
#   - Use during development for quick iteration
build:
	cargo build

# build-release: Compiles the project in release mode.
#   - Slower compilation, optimized runtime performance
#   - Strips debug symbols, smaller binary size
#   - Applies all optimizations (inlining, loop unrolling, etc.)
#   - Output binary: target/release/greeks_test
#   - Use for production deployments and benchmarking
build-release:
	cargo build --release

# check: Analyzes code for errors without producing a binary.
#   - Fastest way to verify code compiles
#   - Skips code generation and linking steps
#   - Useful for quick syntax and type checking during development
#   - Does not produce any output files
check:
	cargo check

# ==============================================================================
# RUN TARGETS
# ==============================================================================

# run: Builds (if needed) and runs the debug binary.
#   - Automatically rebuilds if source files changed
#   - Runs with debug symbols and no optimizations
#   - Pass arguments via ARGS variable: make run ARGS='--ticker TSLA'
#   - The '--' separates cargo args from binary args
run:
	cargo run --bin greeks_test $(if $(ARGS),-- $(ARGS),)

# run-release: Builds (if needed) and runs the release binary.
#   - Automatically rebuilds if source files changed
#   - Runs with full optimizations enabled
#   - Use for performance testing and production-like runs
#   - Pass arguments via ARGS variable: make run-release ARGS='--limit 50'
run-release:
	cargo run --release --bin greeks_test $(if $(ARGS),-- $(ARGS),)

# example: Runs the basic_usage example from the examples/ directory.
#   - Examples demonstrate library usage for other developers
#   - Located in examples/basic_usage.rs
#   - Imports the library as an external consumer would
#   - Useful for validating the public API is ergonomic
example:
	cargo run --example basic_usage

# ==============================================================================
# TEST & LINT TARGETS
# ==============================================================================

# test: Runs all unit and integration tests.
#   - Executes tests in src/ (unit tests) and tests/ (integration tests)
#   - Tests run in parallel by default
#   - Fails if any test fails
#   - Use `cargo test -- --nocapture` to see println! output
test:
	cargo test

# fmt: Formats all Rust source files using rustfmt.
#   - Applies consistent code style (indentation, spacing, etc.)
#   - Configured via rustfmt.toml if present
#   - Modifies files in place
#   - Run before committing to ensure consistent style
fmt:
	cargo fmt

# lint: Runs clippy, Rust's official linter.
#   - Catches common mistakes and suggests improvements
#   - -D warnings treats all warnings as errors (fails the build)
#   - Checks for performance issues, code smells, and anti-patterns
#   - More thorough than `cargo check`
lint:
	cargo clippy -- -D warnings

# ==============================================================================
# CLEANUP TARGETS
# ==============================================================================

# clean: Removes all build artifacts.
#   - Deletes the entire target/ directory
#   - Frees disk space (target/ can grow to several GB)
#   - Next build will recompile everything from scratch
#   - Use when switching branches or troubleshooting build issues
clean:
	cargo clean

# ==============================================================================
# GO + FRONTEND TARGETS
# ==============================================================================

# run_go_front: Starts Go backend and Next.js frontend in topological order.
#   - Starts Go backend first (listens on :8080)
#   - Waits 2 seconds for backend to initialize
#   - Starts Next.js frontend (listens on :3000)
#   - Ctrl+C stops both processes gracefully
#   - View at http://localhost:3000
run_go_front:
	@echo "Starting Go backend on :8080..."
	@cd backend-go && make run & GO_PID=$$!; \
	sleep 2; \
	echo "Starting Next.js frontend on :3000..."; \
	echo "Press Ctrl+C to stop both services"; \
	cd frontend && npm run dev; \
	kill $$GO_PID 2>/dev/null || true

# ==============================================================================
# DOCKER TARGETS
# ==============================================================================

# docker-build: Build Docker images for all services
#   - Builds backend and frontend images
#   - Uses docker-compose build with no cache
docker-build:
	@echo "Building Docker images..."
	docker-compose build --no-cache

# docker-up: Start all services with Docker Compose
#   - Builds images if needed
#   - Starts services in detached mode
#   - Backend: http://localhost:8080
#   - Frontend: http://localhost:3000
#   - View logs with: make docker-logs
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d --build
	@echo ""
	@echo "✓ Services started!"
	@echo "  Backend:  http://localhost:8080"
	@echo "  Frontend: http://localhost:3000"
	@echo ""
	@echo "View logs with: make docker-logs"
	@echo "Stop services with: make docker-down"

# docker-down: Stop and remove all containers
#   - Stops all running containers
#   - Removes containers and networks
#   - Preserves volumes and images
docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

# docker-logs: View logs from all services
#   - Shows logs from backend and frontend
#   - Follows logs in real-time
#   - Press Ctrl+C to exit
docker-logs:
	@echo "Viewing Docker logs (Ctrl+C to exit)..."
	docker-compose logs -f --tail=100

# docker-restart: Restart all services
#   - Stops and starts all containers
#   - Useful after code changes
docker-restart:
	@echo "Restarting Docker services..."
	docker-compose restart
	@echo "✓ Services restarted"
