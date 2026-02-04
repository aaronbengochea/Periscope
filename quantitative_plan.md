# Periscope - Quantitative Analytics Plan

This document outlines the implementation roadmap for Periscope's quantitative analytics components. For detailed technical methodology (formulas, algorithms, implementation details), see [`goal.md`](goal.md).

---

## Overview

The quantitative analytics pipeline transforms raw options chain data into actionable trading signals through a series of analytical stages:

```
Raw Data → IV Extraction → Skew/Smile → Surface → Heston → Relative Value → Signals
```

Each stage builds on the previous, creating a dependency chain that must be implemented in order.

---

## Dependency Graph

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        QUANTITATIVE BUILD ORDER                              │
└─────────────────────────────────────────────────────────────────────────────┘

Phase Q0: IV Extraction
    │
    ├──► Phase Q1: Skew Engine ──────────┐
    │    (delta-space analysis)          │
    │                                    │
    └──► Phase Q2: Smile Engine ─────────┼──► Phase Q3: Surface Builder
         (strike-space analysis,         │    (SSVI parameterization)
          SVI fitting)                   │           │
                                         │           │
                                         └───────────┼──► Phase Q4: Heston Calibration
                                                     │    (stochastic vol model)
                                                     │           │
                                                     │           ▼
                                                     │    Phase Q5: Relative Value
                                                     │    (model vs market comparison)
                                                     │           │
                                                     └───────────┴──► Phase Q6: Trade Signals
                                                                      (systematic generation)

Parallel Development:
  ├── Phase Q7: Risk-Neutral Density (after Q3)
  ├── Phase Q8: Backtesting Framework (after Q6)
  └── Phase Q9: Risk Management (after Q6)
```

---

## Phase Q0: Implied Volatility Extraction

**Goal**: Extract clean implied volatilities from market option prices.

**Depends on**: Options chain data from Massive API (Phase 0 Infrastructure)

### Deliverables

| Component | Description |
|-----------|-------------|
| `IVSolver` | Newton-Raphson implied volatility solver |
| `IVFilter` | Bid-ask spread filtering, sanity bounds |
| `Moneyness` | Log-moneyness and delta calculations |

### Implementation Tasks

- [ ] Implement Black-Scholes pricing formula
- [ ] Implement Newton-Raphson IV solver (convergence < 1e-8)
- [ ] Add initial guess heuristic for ATM options
- [ ] Implement sanity bounds (0.01 < σ < 5.0)
- [ ] Calculate log-moneyness: k = ln(K/F)
- [ ] Filter options with spread > 20% of mid
- [ ] Implement liquidity weighting: w_i = 1/spread_i

### Acceptance Criteria

- [ ] IV extraction completes in < 1ms per option
- [ ] 95%+ convergence rate on valid options
- [ ] Filtered IVs pass sanity checks (no negative, no extreme values)

### Technical Reference

See `goal.md` Section 1: "Implied Volatility Extraction"

---

## Phase Q1: Skew Calculation Engine

**Goal**: Calculate volatility skew across all deltas and expirations.

**Depends on**: Phase Q0 (IV Extraction)

### Deliverables

| Component | Description |
|-----------|-------------|
| `DeltaCalculator` | Black-Scholes delta for calls/puts |
| `DeltaStrikeMapper` | Find strikes at target deltas via root-finding |
| `SkewCalculator` | Risk reversal calculation at standard deltas |
| `SkewTermStructure` | Skew across expirations |

### Implementation Tasks

- [ ] Implement Black-Scholes delta (spot delta convention)
- [ ] Implement delta-to-strike solver (Brent's method)
- [ ] Calculate risk reversal: RR_Δ = σ_call,Δ - σ_put,Δ
- [ ] Support standard deltas: 10Δ, 15Δ, 25Δ, 35Δ
- [ ] Build skew term structure: {(T_i, RR_25(T_i))}
- [ ] Compute skew z-score vs 252-day rolling window
- [ ] Identify skew anomalies (z > 2 or z < -2)

### Acceptance Criteria

- [ ] Delta-strike mapping accurate to 4 decimal places
- [ ] Skew calculations match Bloomberg/OptionMetrics benchmarks
- [ ] Term structure visualization renders correctly

### Technical Reference

See `goal.md` Section 2: "Skew Extraction: Delta-Space Analysis"

---

## Phase Q2: Smile Calculation Engine

**Goal**: Calculate volatility smile across strikes for each expiration using SVI.

**Depends on**: Phase Q0 (IV Extraction)

### Deliverables

| Component | Description |
|-----------|-------------|
| `ButterflyCalculator` | 25-delta butterfly spread calculation |
| `SVIFitter` | SVI parameterization fitting |
| `SmileMetrics` | Curvature, asymmetry, wing slopes |
| `ArbitrageChecker` | SVI arbitrage constraint validation |

### Implementation Tasks

- [ ] Calculate butterfly: BF_25 = (σ_call,25Δ + σ_put,25Δ)/2 - σ_ATM
- [ ] Implement SVI parameterization: w(k) = a + b·[ρ·(k-m) + √((k-m)² + σ²)]
- [ ] Build SVI calibration (SLSQP optimizer)
- [ ] Enforce arbitrage-free constraints:
  - a + b·σ·√(1-ρ²) ≥ 0
  - b ≥ 0, |ρ| < 1, σ > 0
- [ ] Extract smile metrics: ATM vol, ATM skew, curvature, wing slopes
- [ ] Calculate smile asymmetry index

### Acceptance Criteria

- [ ] SVI fits with RMSE < 50 bps
- [ ] All fitted parameters satisfy arbitrage constraints
- [ ] Smile metrics match analytical formulas

### Technical Reference

See `goal.md` Section 3: "Smile Extraction: Strike-Space Analysis"

---

## Phase Q3: Volatility Surface Construction

**Goal**: Build complete 3D arbitrage-free IV surface using SSVI.

**Depends on**: Phase Q1 (Skew Engine), Phase Q2 (Smile Engine)

### Deliverables

| Component | Description |
|-----------|-------------|
| `SurfaceBuilder` | SSVI surface construction |
| `CalendarArbChecker` | Calendar spread arbitrage detection |
| `ButterflyArbChecker` | Butterfly arbitrage (Durrleman condition) |
| `SurfaceInterpolator` | Dense grid interpolation |

### Implementation Tasks

- [ ] Implement SSVI: w(k,T) = θ_T/2 · {1 + ρ·φ(θ_T)·k + √[(φ·k + ρ)² + (1-ρ²)]}
- [ ] Fit global parameters: {ρ, η, γ}
- [ ] Fit per-expiration θ_T from SVI fits
- [ ] Enforce calendar arbitrage: ∂w(k,T)/∂T ≥ 0
- [ ] Enforce butterfly arbitrage: g(k) ≥ 0 (Durrleman condition)
- [ ] Interpolate to dense grid: strikes 0.5·S to 2.0·S, expiries 7d to 2y
- [ ] Implement surface quality metrics

### Acceptance Criteria

- [ ] Surface passes all arbitrage checks
- [ ] Interpolated values smooth and consistent
- [ ] Surface construction completes in < 500ms

### Technical Reference

See `goal.md` Section 4: "Volatility Surface Construction"

---

## Phase Q4: Heston Model Calibration

**Goal**: Fit Heston stochastic volatility model to observed surface.

**Depends on**: Phase Q3 (Surface Builder)

### Deliverables

| Component | Description |
|-----------|-------------|
| `HestonPricer` | Option pricing via characteristic function |
| `HestonCalibrator` | Global optimization calibration |
| `CharacteristicFunction` | Heston characteristic function (Albrecher formulation) |
| `CalibrationQuality` | RMSE, max error, parameter stability metrics |

### Implementation Tasks

- [ ] Implement Heston characteristic function (avoid "Little Heston Trap")
- [ ] Implement Fourier inversion (Gauss-Laguerre quadrature, 64 points)
- [ ] Build call/put pricer via Carr-Madan formula
- [ ] Implement differential evolution optimizer
- [ ] Set parameter bounds: v0, θ ∈ [0.001, 1.0], κ ∈ [0.01, 10.0], σ ∈ [0.01, 2.0], ρ ∈ [-0.99, 0.0]
- [ ] Enforce Feller condition: 2κθ > σ²
- [ ] Implement calibration objective (IV error or price error)
- [ ] Add regularization using prior day's parameters
- [ ] Track daily parameter evolution

### Acceptance Criteria

- [ ] Calibration RMSE < 50 bps
- [ ] Max error < 2%
- [ ] Calibration completes in < 5 seconds
- [ ] Parameters stable day-over-day (no large jumps without regime change)

### Technical Reference

See `goal.md` Section 5: "Heston Model Calibration"

---

## Phase Q5: Relative Value Derivation

**Goal**: Identify rich and cheap options by comparing model vs market IVs.

**Depends on**: Phase Q4 (Heston Calibration)

### Deliverables

| Component | Description |
|-----------|-------------|
| `ModelIVExtractor` | Extract model-implied IV from Heston prices |
| `MispricingCalculator` | Market IV - Model IV calculation |
| `RelativeValueMatrix` | Mispricing by strike × expiration |
| `SignificanceFilter` | Z-score filtering for tradeable signals |

### Implementation Tasks

- [ ] Price each option using calibrated Heston
- [ ] Invert to get model-implied IV: σ_model
- [ ] Calculate mispricing: Δσ = σ_mkt - σ_model
- [ ] Calculate relative mispricing: Δσ_rel = (σ_mkt - σ_model) / σ_model
- [ ] Compute z-score: z_i = Δσ_i / σ_calibration_error
- [ ] Adjust for transaction costs: Δσ_adjusted = Δσ - (spread/2)/vega
- [ ] Build relative value matrix (strike × expiration)
- [ ] Rank options by absolute mispricing

### Acceptance Criteria

- [ ] Mispricing calculations match manual verification
- [ ] Z-score filtering correctly identifies |z| > 2 options
- [ ] Liquidity-adjusted mispricing accounts for spreads

### Technical Reference

See `goal.md` Section 6: "Model-Implied Volatility and Relative Value"

---

## Phase Q6: Trade Signal Generation

**Goal**: Systematic identification of skew/smile trades.

**Depends on**: Phase Q5 (Relative Value)

### Deliverables

| Component | Description |
|-----------|-------------|
| `RiskReversalSignal` | Skew trade signal generation |
| `RatioSpreadSignal` | Gamma picking trade signals |
| `ButterflySignal` | Smile curvature trade signals |
| `CalendarSignal` | Term structure trade signals |
| `KellyCalculator` | Position sizing via Kelly criterion |

### Implementation Tasks

- [ ] Define trade templates with entry/exit rules
- [ ] Implement risk reversal signal:
  - Entry: RR z-score > 2, put Δσ > +2%, call Δσ < -1%
  - Exit: RR z-score < 0.5 or stop at z > 3.5
- [ ] Implement ratio spread signal:
  - Entry: Put Δσ > +3%, call Δσ < -2%, funding check
  - Construction: Sell 1x put, buy N×calls (N=3-4)
- [ ] Implement butterfly signal:
  - Entry: BF z-score > 2, wing Δσ > +2%, ATM |Δσ| < 0.5%
- [ ] Calculate position Greeks (net Δ, Γ, V, Θ)
- [ ] Estimate P&L scenarios (vol up/down, spot up/down)
- [ ] Implement Kelly sizing: f* = (p·b - q) / b
- [ ] Generate alerts with full rationale

### Acceptance Criteria

- [ ] Signals match manual analysis on historical data
- [ ] Greeks calculations accurate
- [ ] Kelly sizing produces reasonable position sizes

### Technical Reference

See `goal.md` Section 9: "Trade Identification Algorithm"

---

## Phase Q7: Risk-Neutral Density Extraction

**Goal**: Extract market-implied probability distributions.

**Depends on**: Phase Q3 (Surface Builder)

### Deliverables

| Component | Description |
|-----------|-------------|
| `BreedenLitzenberger` | RND extraction from option prices |
| `DensityComparator` | Market vs model density comparison |
| `TailFitter` | GPD fit for extreme tails |

### Implementation Tasks

- [ ] Implement Breeden-Litzenberger: f(K) = e^(rT) · ∂²C/∂K²
- [ ] Numerical second derivative via finite differences
- [ ] Interpolate call prices across strikes (cubic spline)
- [ ] Fit Generalized Pareto Distribution to tails
- [ ] Extract Heston density via inverse Fourier transform
- [ ] Compare market RND vs model RND
- [ ] Identify tail divergences (tail risk mispricing)

### Acceptance Criteria

- [ ] Extracted density integrates to 1.0 (±1%)
- [ ] Tail fits smooth and continuous
- [ ] Divergence metrics identify known mispricing

### Technical Reference

See `goal.md` Section 7: "Risk-Neutral Density Extraction (Breeden-Litzenberger)"

---

## Phase Q8: Backtesting Framework

**Goal**: Validate strategy performance on historical data.

**Depends on**: Phase Q6 (Trade Signals)

### Deliverables

| Component | Description |
|-----------|-------------|
| `HistoricalCalibrator` | Replay calibration on historical snapshots |
| `TradeSimulator` | Entry/exit simulation with mark-to-market |
| `PerformanceMetrics` | Sharpe, win rate, drawdown, Calmar |
| `BacktestReport` | Comprehensive strategy report |

### Implementation Tasks

- [ ] Load historical options chain snapshots
- [ ] Replay Heston calibration for each date
- [ ] Generate historical signals using same rules
- [ ] Simulate trade execution (mid or aggressive entry)
- [ ] Daily mark-to-market using next day's surface
- [ ] Track P&L, holding periods, Greeks evolution
- [ ] Calculate performance metrics:
  - Win rate: # winners / # trades (target > 55%)
  - Avg win / Avg loss (target > 1.5)
  - Sharpe ratio (target > 1.5)
  - Max drawdown (target < 15%)
  - Calmar ratio (target > 1.0)
- [ ] Generate backtest reports with visualizations

### Acceptance Criteria

- [ ] Backtest reproducible across runs
- [ ] P&L attribution matches expected behavior
- [ ] Strategy meets performance targets on in-sample data

### Technical Reference

See `goal.md` Section 11: "Backtesting Framework"

---

## Phase Q9: Risk Management

**Goal**: Portfolio-level risk monitoring and limits.

**Depends on**: Phase Q6 (Trade Signals)

### Deliverables

| Component | Description |
|-----------|-------------|
| `GreeksAggregator` | Portfolio-level Greeks |
| `PositionLimits` | Per-name and portfolio limits |
| `ScenarioAnalyzer` | P&L under stress scenarios |
| `RiskAlerts` | Threshold breach notifications |

### Implementation Tasks

- [ ] Aggregate position Greeks:
  - Portfolio Δ = Σ position_i × Δ_i
  - Portfolio Γ = Σ position_i × Γ_i
  - Portfolio V = Σ position_i × V_i
  - Portfolio Θ = Σ position_i × Θ_i
- [ ] Implement position limits:
  - Single name: max 5% of capital
  - Net delta: < 20% of portfolio
  - Net vega: < 2% of portfolio
- [ ] Build scenario analyzer:
  - Spot ±5%, ±10%, ±20%
  - IV ±5, ±10 vol points
  - Combined crash: spot -10%, IV +10
- [ ] Implement real-time alerts when limits breached
- [ ] Track risk metrics evolution over time

### Acceptance Criteria

- [ ] Greeks aggregation matches position-level sum
- [ ] Limits enforced correctly
- [ ] Scenario P&L matches analytical expectations

### Technical Reference

See `goal.md` Section 12: "Risk Management"

---

## Timeline Overview

```
                        Week
Component               1   2   3   4   5   6   7   8   9  10  11  12
─────────────────────────────────────────────────────────────────────
Q0: IV Extraction      [===]
Q1: Skew Engine            [===]
Q2: Smile Engine           [===]
Q3: Surface Builder            [=====]
Q4: Heston Calibration             [=======]
Q5: Relative Value                     [===]
Q6: Trade Signals                          [=====]
Q7: RND Extraction                 [===]
Q8: Backtesting                                [=======]
Q9: Risk Management                                [=====]
─────────────────────────────────────────────────────────────────────
```

**Notes**:
- Q1 and Q2 can be developed in parallel
- Q7 can start after Q3 (does not block Q4-Q6)
- Q8 and Q9 can proceed in parallel after Q6

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| IV extraction latency | < 1ms/option | Benchmark on 1000 options |
| Surface update latency | < 100ms | End-to-end surface refresh |
| Heston calibration time | < 5 seconds | Full calibration cycle |
| Calibration RMSE | < 50 bps | Daily calibration quality |
| Backtest Sharpe ratio | > 1.5 | 3-year backtest |
| Backtest win rate | > 55% | Signal accuracy |
| Backtest max drawdown | < 15% | Risk-adjusted returns |
| Average trade edge | > 2% IV points | Mean mispricing captured |

---

## Integration with Infrastructure

This quantitative plan integrates with [`infrastructure_plan.md`](infrastructure_plan.md) as follows:

| Quant Component | Infrastructure Dependency |
|-----------------|---------------------------|
| Q0: IV Extraction | Phase 0: Options chain data |
| Q1-Q6: Analytics | Phase 1-3: Core services |
| Surface storage | Phase 2: TimescaleDB |
| Real-time updates | Phase 3: WebSocket API |
| Visualization | Phase 4: Frontend surface component |
| Alerting | Phase 5: Notification service |
| Backtesting | Phase 6: Historical data pipeline |

---

## References

- [Gatheral & Jacquier - Arbitrage-free SVI volatility surfaces](https://arxiv.org/abs/1204.0646)
- [Heston Model Analysis and Implementation](https://arxiv.org/pdf/1502.02963)
- [Full and Fast Calibration of the Heston Model](https://eprints.lse.ac.uk/83754/1/Germano_Full%20and%20fast%20calibration_2017.pdf)
- [Breeden-Litzenberger Risk-Neutral Density Extraction](https://www.newyorkfed.org/medialibrary/media/research/staff_reports/sr677.pdf)
