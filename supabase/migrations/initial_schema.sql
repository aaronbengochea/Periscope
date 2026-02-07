-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
CREATE EXTENSION IF NOT EXISTS pg_cron;
CREATE EXTENSION IF NOT EXISTS pg_partman;

-- Options contracts table (master data)
CREATE TABLE IF NOT EXISTS options_contracts (
  id BIGSERIAL PRIMARY KEY,
  ticker TEXT NOT NULL UNIQUE,
  contract_type TEXT NOT NULL CHECK (contract_type IN ('call', 'put')),
  strike_price NUMERIC(10, 2) NOT NULL,
  expiration_date DATE NOT NULL,
  underlying_ticker TEXT NOT NULL,
  exercise_style TEXT CHECK (exercise_style IN ('american', 'european', 'bermudan')),
  shares_per_contract INTEGER DEFAULT 100,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_options_contracts_ticker ON options_contracts(ticker);
CREATE INDEX IF NOT EXISTS idx_options_contracts_underlying ON options_contracts(underlying_ticker);
CREATE INDEX IF NOT EXISTS idx_options_contracts_expiration ON options_contracts(expiration_date);
CREATE INDEX IF NOT EXISTS idx_options_contracts_type ON options_contracts(contract_type);

-- Add comment
COMMENT ON TABLE options_contracts IS 'Master table for options contract metadata';

-- Options quotes table (time-series data) - partitioned by timestamp
CREATE TABLE IF NOT EXISTS options_quotes (
  id BIGSERIAL NOT NULL,
  ticker TEXT NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL,
  bid NUMERIC(10, 4),
  ask NUMERIC(10, 4),
  last_price NUMERIC(10, 4),
  volume BIGINT,
  open_interest BIGINT,
  implied_volatility NUMERIC(10, 6),
  delta NUMERIC(10, 6),
  gamma NUMERIC(10, 6),
  theta NUMERIC(10, 6),
  vega NUMERIC(10, 6),
  rho NUMERIC(10, 6),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);

-- Create default partition
CREATE TABLE IF NOT EXISTS options_quotes_default PARTITION OF options_quotes DEFAULT;

-- Create indexes on partitioned table
CREATE INDEX IF NOT EXISTS idx_options_quotes_ticker ON options_quotes(ticker);
CREATE INDEX IF NOT EXISTS idx_options_quotes_timestamp ON options_quotes(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_options_quotes_ticker_timestamp ON options_quotes(ticker, timestamp DESC);

-- Add comment
COMMENT ON TABLE options_quotes IS 'Time-series partitioned table for historical options quotes and Greeks';
