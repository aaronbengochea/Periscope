# Periscope - Infrastructure Plan

## Vision

Build a production-grade, cloud-native options analytics platform with:
- **Robustness**: Fault-tolerant architecture with graceful degradation
- **Scalability**: Horizontal scaling for compute-intensive analytics
- **Collaboration**: Clear ownership boundaries and well-defined interfaces

---

## Design Principles

1. **Start Small, Iterate Fast**: Build minimal viable components first, validate, then extend
2. **Interface-First Design**: Define contracts before implementation
3. **Test at Every Layer**: Unit → Integration → E2E → Load
4. **Infrastructure as Code**: All cloud resources defined in CDK
5. **Observability from Day One**: Logging, metrics, tracing built-in
6. **Fail Fast, Recover Faster**: Circuit breakers, retries, health checks

---

## System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              FRONTEND (React/Next.js)                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐      │
│  │ Dashboard │  │ Surface  │  │ Trade    │  │ Position │  │ Settings │      │
│  │   View   │  │   3D     │  │ Signals  │  │ Manager  │  │  Panel   │      │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘      │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼ WebSocket / REST
┌─────────────────────────────────────────────────────────────────────────────┐
│                              API GATEWAY (AWS)                               │
│                    Rate Limiting │ Auth │ Routing                           │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
        ┌────────────────────────────┼────────────────────────────┐
        ▼                            ▼                            ▼
┌───────────────┐          ┌───────────────┐          ┌───────────────┐
│  REST API     │          │  WebSocket    │          │  Agent        │
│  Service      │          │  Service      │          │  Service      │
│  (Rust/Axum)  │          │  (Rust/Tokio) │          │  (Python/LG)  │
└───────────────┘          └───────────────┘          └───────────────┘
        │                            │                            │
        └────────────────────────────┼────────────────────────────┘
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CORE ANALYTICS ENGINE (Rust)                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐      │
│  │  Skew    │  │  Smile   │  │ Surface  │  │  Heston  │  │ RelValue │      │
│  │ Engine   │  │ Engine   │  │ Builder  │  │Calibrator│  │  Engine  │      │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘      │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
        ┌────────────────────────────┼────────────────────────────┐
        ▼                            ▼                            ▼
┌───────────────┐          ┌───────────────┐          ┌───────────────┐
│  Market Data  │          │   TimeSeries  │          │    Cache      │
│  Client       │          │   Database    │          │   (Redis)     │
│  (Massive)    │          │  (TimescaleDB)│          │               │
└───────────────┘          └───────────────┘          └───────────────┘
```

---

## Build Order (Dependency Graph)

We build components in topological order, ensuring each layer is testable before moving up.

```
Phase 0: Foundation
    └── Local Development Environment
    └── CI/CD Pipeline Skeleton
    └── Logging & Error Framework

Phase 1: Data Layer
    └── Market Data Client (Massive API)
    └── Data Models & Validation
    └── Local Caching (In-Memory)
    └── Unit Tests for Data Layer

Phase 2: Core Analytics (Rust Library)
    └── Black-Scholes Pricer
    └── IV Solver (Newton-Raphson)
    └── Skew Calculator
    └── Smile Calculator
    └── Surface Builder (SVI)
    └── Integration Tests for Analytics

Phase 3: API Layer
    └── REST API (Axum)
    └── WebSocket Server
    └── API Tests

Phase 4: Persistence
    └── TimescaleDB Schema
    └── Repository Pattern Implementation
    └── Migration Framework
    └── Database Integration Tests

Phase 5: Cloud Infrastructure (CDK)
    └── VPC & Networking
    └── ECS/Fargate Services
    └── RDS (TimescaleDB)
    └── ElastiCache (Redis)
    └── API Gateway
    └── Infrastructure Tests

Phase 6: Frontend
    └── Component Library
    └── Dashboard Layout
    └── 3D Surface Visualization
    └── Real-time Data Binding
    └── E2E Tests

Phase 7: Advanced Analytics
    └── Heston Calibrator
    └── Relative Value Engine
    └── Trade Signal Generator
    └── Performance Benchmarks

Phase 8: Agent Layer (LangGraph)
    └── Research Agent
    └── Trade Recommendation Agent
    └── Portfolio Analysis Agent
    └── Agent Integration Tests

Phase 9: Production Hardening
    └── Load Testing
    └── Chaos Engineering
    └── Runbooks & Alerts
    └── Security Audit
```

---

## Phase 0: Foundation

### 0.1 Local Development Environment

**Deliverables**:
- Docker Compose for local services (DB, Redis, Mock API)
- VS Code devcontainer configuration
- Environment variable management (.env templates)

**Files**:
```
docker/
├── docker-compose.yml          # Local development stack
├── docker-compose.test.yml     # Test environment
└── Dockerfile.dev              # Development container
```

### 0.2 CI/CD Pipeline

**Platform**: GitHub Actions

**Workflows**:
```yaml
.github/workflows/
├── ci.yml                      # Build, test, lint on PR
├── cd-staging.yml              # Deploy to staging on merge
├── cd-production.yml           # Deploy to prod (manual trigger)
└── security-scan.yml           # Weekly dependency audit
```

**CI Pipeline Stages**:
1. `cargo fmt --check` - Format validation
2. `cargo clippy` - Lint
3. `cargo test` - Unit & integration tests
4. `cargo build --release` - Build artifacts
5. Docker image build & push to ECR

### 0.3 Observability Framework

**Logging** (tracing crate):
```rust
// src/telemetry/mod.rs
pub fn init_logging() {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env())
        .json()  // Structured JSON for CloudWatch
        .init();
}
```

**Metrics** (prometheus crate):
```rust
// Key metrics to track
- api_requests_total{endpoint, method, status}
- api_request_duration_seconds{endpoint}
- analytics_calculation_duration_seconds{engine}
- market_data_fetch_duration_seconds
- websocket_connections_active
- cache_hit_rate
```

**Tracing**:
- Instrument all async functions with `#[instrument]`
- Propagate trace IDs through WebSocket messages

---

## Phase 1: Data Layer

### 1.1 Market Data Client

**Module**: `src/client/`

**Classes/Structs**:

```rust
// src/client/mod.rs
pub trait MarketDataClient: Send + Sync {
    async fn get_options_chain(&self, symbol: &str, params: ChainParams) -> Result<OptionsChain>;
    async fn get_quote(&self, symbol: &str) -> Result<Quote>;
    async fn subscribe(&self, symbols: &[&str]) -> Result<Subscription>;
}

// src/client/massive.rs
pub struct MassiveClient {
    http: reqwest::Client,
    base_url: String,
    api_key: String,
    rate_limiter: RateLimiter,
}

// src/client/mock.rs (for testing)
pub struct MockClient {
    responses: HashMap<String, OptionsChain>,
}
```

**Rate Limiting**:
```rust
pub struct RateLimiter {
    permits: Arc<Semaphore>,
    refill_interval: Duration,
}
```

### 1.2 Data Models

**Module**: `src/models/`

```rust
// src/models/option.rs
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OptionContract {
    pub ticker: String,
    pub underlying: String,
    pub strike: Decimal,
    pub expiration: NaiveDate,
    pub contract_type: ContractType,
    pub bid: Option<Decimal>,
    pub ask: Option<Decimal>,
    pub last: Option<Decimal>,
    pub volume: i64,
    pub open_interest: i64,
    pub implied_volatility: Option<f64>,
    pub greeks: Option<Greeks>,
}

// src/models/greeks.rs
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Greeks {
    pub delta: f64,
    pub gamma: f64,
    pub theta: f64,
    pub vega: f64,
    pub rho: Option<f64>,
}

// src/models/surface.rs
#[derive(Debug, Clone)]
pub struct VolatilitySurface {
    pub underlying: String,
    pub timestamp: DateTime<Utc>,
    pub spot: Decimal,
    pub slices: Vec<VolatilitySlice>,  // One per expiration
}

#[derive(Debug, Clone)]
pub struct VolatilitySlice {
    pub expiration: NaiveDate,
    pub tte: f64,  // Time to expiry in years
    pub points: Vec<VolatilityPoint>,
    pub svi_params: Option<SviParams>,
}

#[derive(Debug, Clone)]
pub struct VolatilityPoint {
    pub strike: Decimal,
    pub moneyness: f64,      // ln(K/F)
    pub delta: f64,
    pub iv_bid: Option<f64>,
    pub iv_ask: Option<f64>,
    pub iv_mid: f64,
    pub iv_model: Option<f64>,  // From Heston
    pub mispricing: Option<f64>,
}
```

### 1.3 Caching Layer

**Module**: `src/cache/`

```rust
// src/cache/mod.rs
pub trait Cache: Send + Sync {
    async fn get<T: DeserializeOwned>(&self, key: &str) -> Option<T>;
    async fn set<T: Serialize>(&self, key: &str, value: &T, ttl: Duration);
    async fn invalidate(&self, key: &str);
}

// src/cache/memory.rs (Phase 1)
pub struct InMemoryCache {
    store: Arc<DashMap<String, (Vec<u8>, Instant)>>,
}

// src/cache/redis.rs (Phase 4)
pub struct RedisCache {
    pool: Pool<RedisConnectionManager>,
}
```

**Cache Keys**:
```
options:chain:{symbol}:{expiry}     TTL: 5s
options:quote:{symbol}              TTL: 1s
surface:{symbol}                    TTL: 30s
heston:params:{symbol}              TTL: 1h
```

---

## Phase 2: Core Analytics Engine

### 2.1 Pricing Module

**Module**: `src/analytics/pricing/`

```rust
// src/analytics/pricing/black_scholes.rs
pub struct BlackScholes;

impl BlackScholes {
    pub fn call_price(s: f64, k: f64, t: f64, r: f64, q: f64, sigma: f64) -> f64;
    pub fn put_price(s: f64, k: f64, t: f64, r: f64, q: f64, sigma: f64) -> f64;
    pub fn delta(s: f64, k: f64, t: f64, r: f64, q: f64, sigma: f64, is_call: bool) -> f64;
    pub fn gamma(s: f64, k: f64, t: f64, r: f64, q: f64, sigma: f64) -> f64;
    pub fn vega(s: f64, k: f64, t: f64, r: f64, q: f64, sigma: f64) -> f64;
    pub fn theta(s: f64, k: f64, t: f64, r: f64, q: f64, sigma: f64, is_call: bool) -> f64;
}

// src/analytics/pricing/iv_solver.rs
pub struct IvSolver {
    max_iterations: usize,
    tolerance: f64,
}

impl IvSolver {
    pub fn solve(&self, market_price: f64, s: f64, k: f64, t: f64, r: f64, q: f64, is_call: bool) -> Result<f64>;
}
```

### 2.2 Skew Engine

**Module**: `src/analytics/skew/`

```rust
// src/analytics/skew/mod.rs
pub struct SkewEngine;

impl SkewEngine {
    /// Find strike for target delta via interpolation + root finding
    pub fn find_strike_for_delta(
        slice: &VolatilitySlice,
        target_delta: f64,
        is_call: bool,
    ) -> Result<(Decimal, f64)>;  // (strike, iv)

    /// Calculate risk reversal at given delta
    pub fn risk_reversal(
        slice: &VolatilitySlice,
        delta: f64,  // e.g., 0.25
    ) -> Result<f64>;  // IV_call - IV_put

    /// Calculate butterfly at given delta
    pub fn butterfly(
        slice: &VolatilitySlice,
        delta: f64,
    ) -> Result<f64>;  // (IV_call + IV_put)/2 - IV_atm

    /// Full skew term structure
    pub fn skew_term_structure(
        surface: &VolatilitySurface,
        delta: f64,
    ) -> Vec<(f64, f64)>;  // [(tte, RR), ...]
}
```

### 2.3 Smile Engine

**Module**: `src/analytics/smile/`

```rust
// src/analytics/smile/svi.rs
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SviParams {
    pub a: f64,      // Overall variance level
    pub b: f64,      // Wing slope
    pub rho: f64,    // Skew (-1 to 1)
    pub m: f64,      // Horizontal shift
    pub sigma: f64,  // ATM smoothness
}

impl SviParams {
    /// Total variance at log-moneyness k
    pub fn total_variance(&self, k: f64) -> f64 {
        self.a + self.b * (self.rho * (k - self.m)
            + ((k - self.m).powi(2) + self.sigma.powi(2)).sqrt())
    }

    /// Implied volatility at log-moneyness k and time t
    pub fn iv(&self, k: f64, t: f64) -> f64 {
        (self.total_variance(k) / t).sqrt()
    }
}

pub struct SviCalibrator {
    pub optimizer: Optimizer,
}

impl SviCalibrator {
    pub fn fit(&self, slice: &VolatilitySlice) -> Result<SviParams>;
    pub fn is_arbitrage_free(&self, params: &SviParams) -> bool;
}
```

### 2.4 Surface Builder

**Module**: `src/analytics/surface/`

```rust
// src/analytics/surface/builder.rs
pub struct SurfaceBuilder {
    svi_calibrator: SviCalibrator,
    interpolator: Interpolator,
}

impl SurfaceBuilder {
    /// Build surface from raw options chain
    pub fn build(&self, chain: &OptionsChain, spot: Decimal) -> Result<VolatilitySurface>;

    /// Interpolate IV at arbitrary (strike, expiry)
    pub fn interpolate(&self, surface: &VolatilitySurface, strike: Decimal, tte: f64) -> Result<f64>;

    /// Check for calendar spread arbitrage
    pub fn check_calendar_arbitrage(&self, surface: &VolatilitySurface) -> Vec<ArbitrageViolation>;

    /// Check for butterfly arbitrage
    pub fn check_butterfly_arbitrage(&self, surface: &VolatilitySurface) -> Vec<ArbitrageViolation>;
}
```

### 2.5 Heston Calibrator

**Module**: `src/analytics/heston/`

```rust
// src/analytics/heston/model.rs
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HestonParams {
    pub v0: f64,     // Initial variance
    pub theta: f64,  // Long-term variance
    pub kappa: f64,  // Mean reversion speed
    pub sigma: f64,  // Vol of vol
    pub rho: f64,    // Spot-vol correlation
}

impl HestonParams {
    pub fn feller_satisfied(&self) -> bool {
        2.0 * self.kappa * self.theta > self.sigma.powi(2)
    }
}

// src/analytics/heston/pricer.rs
pub struct HestonPricer;

impl HestonPricer {
    /// Price European call via characteristic function
    pub fn call_price(
        params: &HestonParams,
        s: f64, k: f64, t: f64, r: f64, q: f64,
    ) -> f64;

    /// Implied volatility from Heston price
    pub fn implied_vol(
        params: &HestonParams,
        s: f64, k: f64, t: f64, r: f64, q: f64,
    ) -> f64;
}

// src/analytics/heston/calibrator.rs
pub struct HestonCalibrator {
    optimizer: DifferentialEvolution,
    max_iterations: usize,
}

impl HestonCalibrator {
    pub fn calibrate(
        &self,
        surface: &VolatilitySurface,
        spot: f64,
        r: f64,
        q: f64,
    ) -> Result<CalibrationResult>;
}

pub struct CalibrationResult {
    pub params: HestonParams,
    pub rmse: f64,
    pub max_error: f64,
    pub iterations: usize,
}
```

### 2.6 Relative Value Engine

**Module**: `src/analytics/relvalue/`

```rust
// src/analytics/relvalue/mod.rs
pub struct RelativeValueEngine {
    heston_pricer: HestonPricer,
    iv_solver: IvSolver,
}

impl RelativeValueEngine {
    /// Calculate mispricing for entire surface
    pub fn compute_mispricing(
        &self,
        surface: &VolatilitySurface,
        heston_params: &HestonParams,
        spot: f64, r: f64, q: f64,
    ) -> MispricingSurface;

    /// Find most mispriced options
    pub fn top_signals(
        &self,
        mispricing: &MispricingSurface,
        n: usize,
        min_liquidity: f64,
    ) -> Vec<TradeSignal>;
}

#[derive(Debug, Clone)]
pub struct TradeSignal {
    pub option: OptionContract,
    pub market_iv: f64,
    pub model_iv: f64,
    pub mispricing: f64,        // market - model
    pub mispricing_zscore: f64,
    pub direction: SignalDirection,  // Buy or Sell
    pub confidence: f64,
}
```

---

## Phase 3: API Layer

### 3.1 REST API (Axum)

**Module**: `src/api/`

**Endpoints**:

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/v1/chain/{symbol}` | Options chain |
| GET | `/v1/surface/{symbol}` | Volatility surface |
| GET | `/v1/skew/{symbol}` | Skew metrics |
| GET | `/v1/signals/{symbol}` | Trade signals |
| GET | `/v1/heston/{symbol}` | Calibrated Heston params |
| POST | `/v1/price` | Price arbitrary option |

**Router Structure**:
```rust
// src/api/router.rs
pub fn create_router(state: AppState) -> Router {
    Router::new()
        .route("/health", get(health_check))
        .nest("/v1", v1_routes())
        .layer(TraceLayer::new_for_http())
        .layer(CorsLayer::permissive())
        .with_state(state)
}

fn v1_routes() -> Router<AppState> {
    Router::new()
        .route("/chain/:symbol", get(handlers::get_chain))
        .route("/surface/:symbol", get(handlers::get_surface))
        .route("/skew/:symbol", get(handlers::get_skew))
        .route("/signals/:symbol", get(handlers::get_signals))
        .route("/heston/:symbol", get(handlers::get_heston))
        .route("/price", post(handlers::price_option))
}
```

**App State**:
```rust
// src/api/state.rs
#[derive(Clone)]
pub struct AppState {
    pub market_client: Arc<dyn MarketDataClient>,
    pub cache: Arc<dyn Cache>,
    pub surface_builder: Arc<SurfaceBuilder>,
    pub heston_calibrator: Arc<HestonCalibrator>,
    pub relvalue_engine: Arc<RelativeValueEngine>,
    pub db: Option<Pool<Postgres>>,
}
```

### 3.2 WebSocket Server

**Module**: `src/ws/`

```rust
// src/ws/server.rs
pub struct WebSocketServer {
    subscriptions: Arc<DashMap<String, Vec<Sender<Message>>>>,
}

impl WebSocketServer {
    pub async fn handle_connection(&self, socket: WebSocket);
    pub async fn broadcast(&self, symbol: &str, update: SurfaceUpdate);
}

// Message types
#[derive(Serialize, Deserialize)]
#[serde(tag = "type")]
pub enum WsMessage {
    Subscribe { symbols: Vec<String> },
    Unsubscribe { symbols: Vec<String> },
    SurfaceUpdate { symbol: String, surface: VolatilitySurface },
    SignalAlert { signal: TradeSignal },
    Error { message: String },
}
```

---

## Phase 4: Persistence Layer

### 4.1 Database Schema (TimescaleDB)

```sql
-- Hypertable for time-series options data
CREATE TABLE options_snapshots (
    time TIMESTAMPTZ NOT NULL,
    symbol TEXT NOT NULL,
    underlying TEXT NOT NULL,
    strike DECIMAL NOT NULL,
    expiration DATE NOT NULL,
    contract_type TEXT NOT NULL,
    bid DECIMAL,
    ask DECIMAL,
    iv_mid DOUBLE PRECISION,
    delta DOUBLE PRECISION,
    gamma DOUBLE PRECISION,
    theta DOUBLE PRECISION,
    vega DOUBLE PRECISION,
    volume BIGINT,
    open_interest BIGINT
);
SELECT create_hypertable('options_snapshots', 'time');
CREATE INDEX idx_options_symbol_time ON options_snapshots (symbol, time DESC);

-- Calibrated Heston parameters
CREATE TABLE heston_calibrations (
    time TIMESTAMPTZ NOT NULL,
    symbol TEXT NOT NULL,
    v0 DOUBLE PRECISION,
    theta DOUBLE PRECISION,
    kappa DOUBLE PRECISION,
    sigma DOUBLE PRECISION,
    rho DOUBLE PRECISION,
    rmse DOUBLE PRECISION,
    PRIMARY KEY (symbol, time)
);
SELECT create_hypertable('heston_calibrations', 'time');

-- Trade signals
CREATE TABLE trade_signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    time TIMESTAMPTZ NOT NULL,
    symbol TEXT NOT NULL,
    strike DECIMAL NOT NULL,
    expiration DATE NOT NULL,
    contract_type TEXT NOT NULL,
    direction TEXT NOT NULL,
    market_iv DOUBLE PRECISION,
    model_iv DOUBLE PRECISION,
    mispricing DOUBLE PRECISION,
    confidence DOUBLE PRECISION,
    status TEXT DEFAULT 'pending'
);
CREATE INDEX idx_signals_time ON trade_signals (time DESC);
```

### 4.2 Repository Pattern

```rust
// src/repository/mod.rs
#[async_trait]
pub trait OptionsRepository: Send + Sync {
    async fn save_snapshot(&self, snapshot: &OptionsSnapshot) -> Result<()>;
    async fn get_history(&self, symbol: &str, from: DateTime<Utc>, to: DateTime<Utc>) -> Result<Vec<OptionsSnapshot>>;
}

#[async_trait]
pub trait HestonRepository: Send + Sync {
    async fn save_calibration(&self, symbol: &str, params: &HestonParams, rmse: f64) -> Result<()>;
    async fn get_latest(&self, symbol: &str) -> Result<Option<HestonParams>>;
    async fn get_history(&self, symbol: &str, days: i32) -> Result<Vec<(DateTime<Utc>, HestonParams)>>;
}

#[async_trait]
pub trait SignalRepository: Send + Sync {
    async fn save_signal(&self, signal: &TradeSignal) -> Result<Uuid>;
    async fn get_pending(&self, symbol: &str) -> Result<Vec<TradeSignal>>;
    async fn update_status(&self, id: Uuid, status: SignalStatus) -> Result<()>;
}
```

---

## Phase 5: AWS Infrastructure (CDK)

### 5.1 CDK Project Structure

```
infra/
├── bin/
│   └── periscope.ts           # CDK app entry point
├── lib/
│   ├── stacks/
│   │   ├── network-stack.ts   # VPC, subnets, security groups
│   │   ├── database-stack.ts  # RDS, ElastiCache
│   │   ├── compute-stack.ts   # ECS, Fargate
│   │   ├── api-stack.ts       # API Gateway, Lambda authorizer
│   │   └── monitoring-stack.ts # CloudWatch, alarms
│   └── constructs/
│       ├── fargate-service.ts
│       └── timescale-db.ts
├── cdk.json
└── package.json
```

### 5.2 Network Stack

```typescript
// lib/stacks/network-stack.ts
export class NetworkStack extends cdk.Stack {
    public readonly vpc: ec2.Vpc;
    public readonly apiSecurityGroup: ec2.SecurityGroup;
    public readonly dbSecurityGroup: ec2.SecurityGroup;

    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        this.vpc = new ec2.Vpc(this, 'PeriscopeVpc', {
            maxAzs: 2,
            natGateways: 1,
            subnetConfiguration: [
                { name: 'Public', subnetType: ec2.SubnetType.PUBLIC },
                { name: 'Private', subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS },
                { name: 'Isolated', subnetType: ec2.SubnetType.PRIVATE_ISOLATED },
            ],
        });

        // Security groups...
    }
}
```

### 5.3 Compute Stack (ECS Fargate)

```typescript
// lib/stacks/compute-stack.ts
export class ComputeStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props: ComputeStackProps) {
        super(scope, id, props);

        const cluster = new ecs.Cluster(this, 'PeriscopeCluster', {
            vpc: props.vpc,
            containerInsights: true,
        });

        // API Service
        const apiService = new ecs_patterns.ApplicationLoadBalancedFargateService(
            this, 'ApiService', {
                cluster,
                cpu: 512,
                memoryLimitMiB: 1024,
                desiredCount: 2,
                taskImageOptions: {
                    image: ecs.ContainerImage.fromEcrRepository(props.apiRepo),
                    containerPort: 8080,
                    environment: {
                        RUST_LOG: 'info',
                        DATABASE_URL: props.databaseUrl,
                        REDIS_URL: props.redisUrl,
                    },
                    secrets: {
                        MASSIVE_API_KEY: ecs.Secret.fromSecretsManager(props.apiKeySecret),
                    },
                },
                publicLoadBalancer: false,
            }
        );

        // Auto-scaling
        const scaling = apiService.service.autoScaleTaskCount({
            minCapacity: 2,
            maxCapacity: 10,
        });
        scaling.scaleOnCpuUtilization('CpuScaling', {
            targetUtilizationPercent: 70,
        });
    }
}
```

### 5.4 Database Stack

```typescript
// lib/stacks/database-stack.ts
export class DatabaseStack extends cdk.Stack {
    public readonly timescaleDb: rds.DatabaseInstance;
    public readonly redis: elasticache.CfnCacheCluster;

    constructor(scope: Construct, id: string, props: DatabaseStackProps) {
        super(scope, id, props);

        // TimescaleDB on RDS PostgreSQL
        this.timescaleDb = new rds.DatabaseInstance(this, 'TimescaleDb', {
            engine: rds.DatabaseInstanceEngine.postgres({
                version: rds.PostgresEngineVersion.VER_15,
            }),
            instanceType: ec2.InstanceType.of(
                ec2.InstanceClass.R6G,
                ec2.InstanceSize.LARGE
            ),
            vpc: props.vpc,
            vpcSubnets: { subnetType: ec2.SubnetType.PRIVATE_ISOLATED },
            allocatedStorage: 100,
            maxAllocatedStorage: 500,
            multiAz: true,
            deletionProtection: true,
        });

        // ElastiCache Redis
        this.redis = new elasticache.CfnCacheCluster(this, 'Redis', {
            cacheNodeType: 'cache.r6g.large',
            engine: 'redis',
            numCacheNodes: 1,
            vpcSecurityGroupIds: [props.cacheSecurityGroup.securityGroupId],
        });
    }
}
```

---

## Phase 6: Frontend

### 6.1 Tech Stack

- **Framework**: Next.js 14 (App Router)
- **UI**: Tailwind CSS + shadcn/ui
- **State**: Zustand
- **3D**: Three.js / React Three Fiber
- **Charts**: Recharts
- **WebSocket**: Native WebSocket + reconnect logic

### 6.2 Component Architecture

```
frontend/
├── app/
│   ├── layout.tsx
│   ├── page.tsx                    # Dashboard
│   ├── surface/[symbol]/page.tsx   # 3D Surface view
│   ├── signals/page.tsx            # Trade signals
│   └── settings/page.tsx
├── components/
│   ├── ui/                         # shadcn components
│   ├── layout/
│   │   ├── Sidebar.tsx
│   │   ├── Header.tsx
│   │   └── Footer.tsx
│   ├── dashboard/
│   │   ├── SymbolSelector.tsx
│   │   ├── SkewCard.tsx
│   │   ├── SmileCard.tsx
│   │   └── QuickStats.tsx
│   ├── surface/
│   │   ├── Surface3D.tsx           # Three.js surface
│   │   ├── SurfaceControls.tsx
│   │   ├── ColorLegend.tsx
│   │   └── SliceView.tsx
│   ├── signals/
│   │   ├── SignalTable.tsx
│   │   ├── SignalCard.tsx
│   │   └── SignalFilters.tsx
│   └── charts/
│       ├── SkewChart.tsx
│       ├── TermStructureChart.tsx
│       └── HestonParamsChart.tsx
├── hooks/
│   ├── useWebSocket.ts
│   ├── useSurface.ts
│   └── useSignals.ts
├── stores/
│   ├── surfaceStore.ts
│   ├── signalStore.ts
│   └── settingsStore.ts
└── lib/
    ├── api.ts                      # REST client
    └── ws.ts                       # WebSocket client
```

### 6.3 3D Surface Component

```tsx
// components/surface/Surface3D.tsx
import { Canvas } from '@react-three/fiber';
import { OrbitControls } from '@react-three/drei';

interface Surface3DProps {
    surface: VolatilitySurface;
    colorMode: 'iv' | 'mispricing';
}

export function Surface3D({ surface, colorMode }: Surface3DProps) {
    const geometry = useMemo(() => buildSurfaceGeometry(surface), [surface]);
    const colors = useMemo(() =>
        colorMode === 'mispricing'
            ? mispricingColors(surface)
            : ivColors(surface),
        [surface, colorMode]
    );

    return (
        <Canvas camera={{ position: [5, 5, 5] }}>
            <ambientLight intensity={0.5} />
            <directionalLight position={[10, 10, 5]} />
            <mesh geometry={geometry}>
                <meshStandardMaterial vertexColors />
            </mesh>
            <OrbitControls />
            <axesHelper args={[5]} />
        </Canvas>
    );
}

function mispricingColors(surface: VolatilitySurface): Float32Array {
    // Map mispricing to color:
    // -3% -> blue (0, 0, 1)
    //  0% -> white (1, 1, 1)
    // +3% -> red (1, 0, 0)
}
```

---

## Phase 7: Agent Layer (LangGraph)

### 7.1 Agent Architecture

```
agents/
├── src/
│   ├── graphs/
│   │   ├── research_agent.py      # Market research
│   │   ├── trade_agent.py         # Trade recommendations
│   │   └── portfolio_agent.py     # Portfolio analysis
│   ├── tools/
│   │   ├── market_data.py         # Fetch from Periscope API
│   │   ├── news_search.py         # News API integration
│   │   ├── sentiment.py           # Sentiment analysis
│   │   └── risk_calculator.py     # Position risk
│   ├── prompts/
│   │   ├── research.py
│   │   ├── trade.py
│   │   └── portfolio.py
│   └── state/
│       └── schemas.py             # Pydantic state schemas
├── requirements.txt
└── Dockerfile
```

### 7.2 Research Agent (LangGraph)

```python
# agents/src/graphs/research_agent.py
from langgraph.graph import StateGraph, END
from langchain_anthropic import ChatAnthropic
from pydantic import BaseModel

class ResearchState(BaseModel):
    symbol: str
    surface_data: Optional[dict] = None
    news: list[str] = []
    sentiment: Optional[float] = None
    analysis: Optional[str] = None
    recommendations: list[str] = []

def fetch_surface(state: ResearchState) -> ResearchState:
    """Fetch volatility surface from Periscope API"""
    response = requests.get(f"{PERISCOPE_API}/v1/surface/{state.symbol}")
    state.surface_data = response.json()
    return state

def fetch_news(state: ResearchState) -> ResearchState:
    """Fetch recent news for symbol"""
    # ... news API call
    return state

def analyze(state: ResearchState) -> ResearchState:
    """LLM analysis of surface + news"""
    llm = ChatAnthropic(model="claude-sonnet-4-20250514")
    prompt = RESEARCH_PROMPT.format(
        symbol=state.symbol,
        surface=state.surface_data,
        news=state.news,
    )
    state.analysis = llm.invoke(prompt).content
    return state

def build_research_graph():
    graph = StateGraph(ResearchState)

    graph.add_node("fetch_surface", fetch_surface)
    graph.add_node("fetch_news", fetch_news)
    graph.add_node("analyze", analyze)

    graph.set_entry_point("fetch_surface")
    graph.add_edge("fetch_surface", "fetch_news")
    graph.add_edge("fetch_news", "analyze")
    graph.add_edge("analyze", END)

    return graph.compile()
```

### 7.3 Trade Recommendation Agent

```python
# agents/src/graphs/trade_agent.py
class TradeState(BaseModel):
    symbol: str
    signals: list[dict] = []
    risk_params: dict = {}
    recommendations: list[TradeRecommendation] = []
    reasoning: str = ""

class TradeRecommendation(BaseModel):
    trade_type: str  # "risk_reversal", "ratio_spread", "butterfly"
    legs: list[TradeLeg]
    max_loss: float
    expected_edge: float
    confidence: float

def evaluate_signals(state: TradeState) -> TradeState:
    """Filter and rank trade signals"""
    # Get signals from Periscope API
    signals = fetch_signals(state.symbol)
    # Filter by liquidity, confidence
    state.signals = [s for s in signals if s['confidence'] > 0.6]
    return state

def construct_trades(state: TradeState) -> TradeState:
    """Construct trade structures from signals"""
    llm = ChatAnthropic(model="claude-sonnet-4-20250514")
    prompt = TRADE_CONSTRUCTION_PROMPT.format(
        signals=state.signals,
        risk_params=state.risk_params,
    )
    # LLM generates trade structures
    state.recommendations = parse_trade_recommendations(llm.invoke(prompt))
    return state

def risk_check(state: TradeState) -> TradeState:
    """Validate risk limits"""
    for rec in state.recommendations:
        if rec.max_loss > state.risk_params.get('max_loss_per_trade', 1000):
            rec.confidence *= 0.5  # Downgrade confidence
    return state
```

---

## Phase 8: Testing Strategy

### 8.1 Test Pyramid

```
                    ┌─────────────┐
                    │    E2E      │  <- Playwright (critical paths)
                    │   Tests     │
                   ┌┴─────────────┴┐
                   │  Integration  │  <- API + DB tests
                   │    Tests      │
                  ┌┴───────────────┴┐
                  │   Unit Tests    │  <- Analytics, models, utils
                  └─────────────────┘
```

### 8.2 Test Organization

```
tests/
├── unit/
│   ├── analytics/
│   │   ├── test_black_scholes.rs
│   │   ├── test_iv_solver.rs
│   │   ├── test_svi.rs
│   │   └── test_heston.rs
│   ├── models/
│   │   └── test_option.rs
│   └── cache/
│       └── test_memory_cache.rs
├── integration/
│   ├── test_api.rs
│   ├── test_websocket.rs
│   └── test_repository.rs
└── e2e/
    └── playwright/
        ├── dashboard.spec.ts
        └── surface.spec.ts
```

### 8.3 Performance Benchmarks

```rust
// benches/analytics_bench.rs
use criterion::{criterion_group, criterion_main, Criterion};

fn bench_black_scholes(c: &mut Criterion) {
    c.bench_function("bs_call_price", |b| {
        b.iter(|| BlackScholes::call_price(100.0, 100.0, 0.25, 0.05, 0.02, 0.20))
    });
}

fn bench_iv_solver(c: &mut Criterion) {
    let solver = IvSolver::default();
    c.bench_function("iv_solve", |b| {
        b.iter(|| solver.solve(5.50, 100.0, 100.0, 0.25, 0.05, 0.02, true))
    });
}

fn bench_heston_price(c: &mut Criterion) {
    let params = HestonParams { v0: 0.04, theta: 0.04, kappa: 2.0, sigma: 0.3, rho: -0.7 };
    c.bench_function("heston_call", |b| {
        b.iter(|| HestonPricer::call_price(&params, 100.0, 100.0, 0.25, 0.05, 0.02))
    });
}

criterion_group!(benches, bench_black_scholes, bench_iv_solver, bench_heston_price);
criterion_main!(benches);
```

**Target Latencies**:
| Operation | Target | Measured |
|-----------|--------|----------|
| BS price | < 100ns | - |
| IV solve | < 10μs | - |
| SVI fit (1 slice) | < 1ms | - |
| Heston price | < 100μs | - |
| Heston calibration | < 5s | - |
| Full surface build | < 100ms | - |

---

## Completed Checklist

### Phase 0: Foundation
- [x] Project initialization and structure
- [x] Rust project setup with modular architecture
- [x] Makefile for build automation
- [x] Environment configuration management
- [x] Error handling infrastructure
- [ ] Docker Compose for local development
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Structured logging setup

### Phase 1: Data Layer
- [x] Massive API client integration
- [x] Options chain snapshot fetching with Greeks
- [x] Data models (OptionContract, Greeks)
- [ ] Rate limiting
- [ ] In-memory caching
- [ ] Mock client for testing

### Phase 2: Core Analytics
- [ ] Black-Scholes pricer
- [ ] IV solver (Newton-Raphson)
- [ ] Skew calculator
- [ ] Smile calculator (SVI)
- [ ] Surface builder
- [ ] Heston calibrator
- [ ] Relative value engine

### Phase 3: API Layer
- [ ] REST API (Axum)
- [ ] WebSocket server
- [ ] API tests

### Phase 4: Persistence
- [ ] TimescaleDB schema
- [ ] Repository implementations
- [ ] Migrations

### Phase 5: Cloud Infrastructure
- [ ] CDK project setup
- [ ] VPC & networking
- [ ] ECS/Fargate services
- [ ] RDS + ElastiCache
- [ ] API Gateway

### Phase 6: Frontend
- [ ] Next.js project setup
- [ ] Component library
- [ ] Dashboard layout
- [ ] 3D surface visualization
- [ ] WebSocket integration

### Phase 7: Advanced Analytics
- [ ] Heston calibration optimization
- [ ] Trade signal generator
- [ ] Backtesting framework

### Phase 8: Agent Layer
- [ ] LangGraph project setup
- [ ] Research agent
- [ ] Trade recommendation agent
- [ ] Portfolio analysis agent

---

## Notes

*Architecture decisions, meeting notes, and blockers go here.*

### ADR-001: Rust for Core Analytics
**Decision**: Use Rust for all performance-critical analytics code.
**Rationale**: Sub-millisecond latency requirements, memory safety, excellent async support.
**Alternatives considered**: C++ (rejected: safety concerns), Python (rejected: too slow for real-time).

### ADR-002: TimescaleDB for Time-Series
**Decision**: Use TimescaleDB (PostgreSQL extension) for options snapshots.
**Rationale**: SQL familiarity, automatic partitioning, compression, continuous aggregates.
**Alternatives considered**: InfluxDB (rejected: less mature), QuestDB (rejected: smaller ecosystem).
