# ==============================================================================
# Periscope Makefile
# ==============================================================================
# This Makefile provides convenient shortcuts for common Rust development tasks.
# Run `make` or `make help` to see all available commands.
# ==============================================================================

.PHONY: build build-release run run-release check test fmt lint clean help

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
