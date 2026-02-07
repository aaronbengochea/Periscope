# Database Migrations

This directory contains SQL migration files for the Periscope database schema.

## Migration Files

- `initial_schema.sql` - Initial database schema with options_contracts and options_quotes tables

## Running Migrations

### Option 1: Supabase CLI (Recommended)

If you have the Supabase CLI installed:

```bash
# Login to Supabase
supabase login

# Link to your project
supabase link --project-ref kkrwstbjorlomjxdrimv

# Run migrations
supabase db push
```

### Option 2: Manual SQL Execution

1. Go to your Supabase dashboard: https://supabase.com/dashboard/project/kkrwstbjorlomjxdrimv
2. Navigate to **SQL Editor**
3. Copy and paste the contents of each migration file in order
4. Execute the SQL

### Option 3: Using psql

```bash
# Set your database connection string
export DATABASE_URL="postgresql://postgres:[YOUR-PASSWORD]@db.kkrwstbjorlomjxdrimv.supabase.co:5432/postgres"

# Run the migration
psql $DATABASE_URL < supabase/migrations/initial_schema.sql
```

## Schema Overview

### Tables

1. **options_contracts** - Master table for options contract metadata
   - Primary key: `id`
   - Unique: `ticker`
   - Indexes: ticker, underlying_ticker, expiration_date, contract_type

2. **options_quotes** - Time-series partitioned table for historical quotes
   - Primary key: `(id, timestamp)`
   - Partitioned by: `timestamp` (RANGE)
   - Indexes: ticker, timestamp, (ticker + timestamp)
   - Columns: bid, ask, last_price, volume, open_interest, Greeks (delta, gamma, theta, vega, rho), implied_volatility

### Extensions

- `pg_stat_statements` - Query performance monitoring
- `pg_cron` - Job scheduler for PostgreSQL
- `pg_partman` - Partition management for time-series data

## Creating New Migrations

When you need to modify the schema:

1. Create a new migration file with a timestamp:
   ```bash
   touch supabase/migrations/$(date +%Y%m%d%H%M%S)_description.sql
   ```

2. Add your SQL changes (ALTER TABLE, CREATE INDEX, etc.)

3. Test the migration on a development database first

4. Commit the migration file to version control

## Rollback

Currently, rollbacks are manual. If you need to undo a migration:

1. Create a new migration that reverses the changes
2. Or manually execute the reverse SQL in the Supabase SQL Editor

## Best Practices

- ✅ Always use `IF NOT EXISTS` and `IF EXISTS` for idempotency
- ✅ Test migrations on a development database first
- ✅ Never modify existing migration files once merged to main
- ✅ Use transactions for complex migrations (BEGIN; ... COMMIT;)
- ✅ Add comments to explain complex schema changes
