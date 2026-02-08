# Periscope - Quantitative Goals

## Big Picture Objective

Build a comprehensive options volatility analytics platform that:
1. Calculates volatility **skew** and **smile** across all strikes and expirations
2. Constructs a complete **3D implied volatility surface**
3. Fits a **Heston stochastic volatility model** to the observed surface
4. Derives **relative value** by comparing model-implied vol to market vol
5. Visualizes mispricing as a **temperature-mapped 3D surface** (blue = cheap, red = rich)
6. Enables **programmatic trading** of identified mispricings

---

## Core Concepts (Refined from idea.md)

### Volatility Skew
**Definition**: The difference in implied volatility between equidistant OTM puts and calls at the same delta.

**Calculation**:
- Take delta-equivalent options (e.g., 25-delta put vs 25-delta call)
- Compare their implied volatilities
- Skew = IV(put) - IV(call)
- Positive skew indicates puts are relatively expensive (typical in equities due to hedging demand)

**What it tells us**: Directional sentiment. High put skew signals fear of downside; the market is pricing crash protection at a premium.

**Example from idea.md**:
- 12-delta put: 62% IV
- 12-delta call: 53% IV
- Skew: +9% → Puts are rich, calls are cheap

### Volatility Smile
**Definition**: The pattern of implied volatility across strikes at a fixed expiration, comparing strikes equidistant from ATM in nominal terms (not delta).

**Calculation**:
- Take strike-equivalent options (e.g., ATM ± $50)
- Compare their implied volatilities
- Smile curvature = how IV changes as you move away from ATM

**What it tells us**: Tail risk pricing. A steep smile indicates the market expects fat tails (large moves more likely than normal distribution suggests).

### Skew vs Smile
| Metric | Comparison Basis | Reveals |
|--------|------------------|---------|
| Skew | Delta-equivalent | Directional bias |
| Smile | Strike-equivalent | Tail probability pricing |

### The Volatility Surface
Combining skew and smile analysis across **all expirations** creates a 3D surface:
- **X-axis**: Strike (or moneyness/delta)
- **Y-axis**: Days to expiration
- **Z-axis**: Implied volatility

This surface represents the market's complete view on volatility across all strikes and time horizons.

### Heston Model
**Why Heston over Black-Scholes?**
- Black-Scholes assumes constant volatility (unrealistic)
- Heston assumes volatility is stochastic (random, mean-reverting)
- Markets empirically display stochastic vol behavior

**Heston Parameters**:
- `v0`: Initial variance
- `θ` (theta): Long-term variance mean
- `κ` (kappa): Mean reversion speed
- `σ` (sigma): Volatility of volatility
- `ρ` (rho): Correlation between spot and vol

**What Heston enables**:
- Simulate forward volatility surfaces
- Estimate tail probabilities
- Extract risk-neutral densities
- Generate model-implied IVs for comparison

### Relative Value Framework
1. Observe market IV for each option
2. Fit Heston model to current surface
3. Calculate model-implied IV for each option
4. Compare: `Mispricing = Market IV - Model IV`
5. Positive = option is **rich** (overpriced)
6. Negative = option is **cheap** (underpriced)

### Kelly Criterion for Sizing
Once edge is identified, size positions using Kelly:
```
f* = (p * b - q) / b
```
Where:
- `f*` = fraction of capital to wager
- `p` = probability of winning
- `b` = odds received (payout ratio)
- `q` = probability of losing (1 - p)

---

## Implementation Phases

### Phase 1: Data Infrastructure
**Goal**: Reliable, low-latency options chain data

- [ ] Fetch complete options chains for target underlyings
- [ ] Store historical snapshots for backtesting
- [ ] Normalize data: strikes, expirations, bid/ask, IV, Greeks
- [ ] Handle corporate actions (splits, dividends)
- [ ] Support multiple underlyings simultaneously

**Deliverables**:
- Options chain ingestion pipeline
- Historical data storage (time-series DB)
- Data quality validation

---

### Phase 2: Skew Calculation Engine
**Goal**: Calculate volatility skew across all deltas and expirations

- [ ] Implement delta calculation (Black-Scholes delta)
- [ ] Build delta-equivalent option pairing logic
- [ ] Calculate skew at standard deltas (10, 25, 50)
- [ ] Compute reversal prices (call - put premium at same delta)
- [ ] Track skew time series for each underlying
- [ ] Identify skew anomalies and mean-reversion opportunities

**Deliverables**:
- `SkewCalculator` service
- Skew term structure (skew vs expiration)
- Historical skew percentile rankings

---

### Phase 3: Smile Calculation Engine
**Goal**: Calculate volatility smile across all strikes for each expiration

- [ ] Implement strike-equivalent pairing (equidistant from ATM)
- [ ] Calculate smile curvature metrics
- [ ] Measure wing steepness (OTM put slope vs OTM call slope)
- [ ] Identify smile asymmetry
- [ ] Track smile evolution over time

**Deliverables**:
- `SmileCalculator` service
- Smile shape metrics (curvature, asymmetry, kurtosis-implied)
- Per-expiration smile snapshots

---

### Phase 4: Volatility Surface Construction
**Goal**: Build complete 3D IV surface from market data

- [ ] Aggregate skew and smile data across all expirations
- [ ] Interpolate missing strikes/expirations (cubic spline or SABR)
- [ ] Handle illiquid options (bid-ask spread filtering)
- [ ] Ensure arbitrage-free surface (calendar spread, butterfly constraints)
- [ ] Smooth surface while preserving market information

**Deliverables**:
- `SurfaceBuilder` service
- Surface data structure (matrix of IVs by strike × expiration)
- Surface quality metrics

---

### Phase 5: Heston Model Calibration
**Goal**: Fit Heston model to observed volatility surface

- [ ] Implement Heston pricing formula (semi-closed form via characteristic function)
- [ ] Build calibration routine (minimize error between model and market prices)
- [ ] Use global optimization (differential evolution, particle swarm)
- [ ] Validate calibration quality (RMSE, parameter stability)
- [ ] Store calibrated parameters daily for time-series analysis
- [ ] Track parameter evolution (κ, θ, σ, ρ trending)

**Deliverables**:
- `HestonCalibrator` service
- Daily calibrated parameters
- Calibration quality dashboard

---

### Phase 6: Relative Value Derivation
**Goal**: Identify rich and cheap options by comparing model vs market

- [ ] Generate model-implied IV for every option from calibrated Heston
- [ ] Calculate mispricing: `Market IV - Model IV`
- [ ] Normalize mispricing (z-score relative to historical distribution)
- [ ] Rank options by relative value
- [ ] Filter by liquidity and transaction cost thresholds
- [ ] Identify trade candidates (rich to sell, cheap to buy)

**Deliverables**:
- `RelativeValueEngine` service
- Mispricing matrix (by strike × expiration)
- Trade signal generation

---

### Phase 7: 3D Visualization
**Goal**: Temperature-mapped relative value surface

**Color Mapping**:
- Deep blue: Significantly cheap (model IV >> market IV)
- Light blue: Moderately cheap
- White/neutral: Fair value
- Orange: Moderately rich
- Deep red/lava: Significantly rich (market IV >> model IV)

**Visualization Features**:
- [ ] Interactive 3D surface plot (WebGL/Three.js)
- [ ] Rotate, zoom, slice by expiration
- [ ] Hover for option details (strike, IV, mispricing, Greeks)
- [ ] Time-lapse animation (surface evolution over days)
- [ ] Overlay current positions
- [ ] Export snapshots for reporting

**Deliverables**:
- Frontend 3D surface component
- Real-time surface updates via WebSocket
- Historical surface replay

---

### Phase 8: Trade Signal Generation
**Goal**: Systematic identification of skew/smile trades

**Trade Types**:
1. **Ratio Spreads**: Sell rich OTM puts, buy multiple cheap OTM calls (3:1, 4:1)
2. **Risk Reversals**: Sell OTM put, buy OTM call (delta-neutral skew trade)
3. **Butterflies**: Exploit smile curvature mispricing
4. **Calendar Spreads**: Exploit term structure mispricing

- [ ] Define trade templates with entry/exit rules
- [ ] Calculate position Greeks (net delta, gamma, vega, theta)
- [ ] Estimate P&L scenarios (vol up/down, spot up/down)
- [ ] Apply Kelly sizing based on edge magnitude
- [ ] Generate trade alerts with full rationale

**Deliverables**:
- `TradeSignalGenerator` service
- Signal dashboard with confidence scores
- Backtesting framework

---

### Phase 9: Programmatic Execution (Future)
**Goal**: Automated trade execution with risk controls

- [ ] Broker API integration (IBKR, Schwab, etc.)
- [ ] Order management system
- [ ] Position monitoring and Greeks tracking
- [ ] Stop-loss and profit-taking automation
- [ ] Portfolio-level risk limits
- [ ] Execution quality analysis (slippage tracking)

**Deliverables**:
- Execution engine
- Risk management dashboard
- Trade journal with P&L attribution

---

## Success Metrics

| Metric | Target |
|--------|--------|
| Surface update latency | < 100ms |
| Heston calibration time | < 5 seconds |
| Backtest Sharpe ratio | > 1.5 |
| Win rate on relative value trades | > 55% |
| Average trade edge (IV points) | > 2% |

---

## Key Insights from idea.md

1. **Structural richness**: Puts are systematically expensive due to institutional hedging demand. This creates persistent alpha opportunity.

2. **Skew as signal**: Elevated skew (put IV >> call IV) often mean-reverts, especially after fear spikes.

3. **Asymmetric payoffs**: Ratio spreads (sell 1 put, buy 3-4 calls) create convex payoffs with defined downside.

4. **Momentum correlation**: Volatility mispricing correlates with price momentum. Coupling vol signals with momentum studies improves timing.

5. **Daily recalibration**: Markets evolve; recalibrate Heston daily to capture regime changes.

---

## Architecture Note

All analytics (skew, smile, surface, Heston, relative value) will be implemented in **Rust** for:
- Sub-millisecond calculation latency
- Memory safety in production
- Concurrent processing of multiple underlyings
- Direct integration with low-latency trading systems

The frontend will consume analytics via WebSocket for real-time visualization updates.

---

## Technical Methodology: Detailed Implementation

This section provides rigorous, step-by-step technical specifications for each analytical component. The methodology draws from established quantitative finance literature and practitioner best practices.

### References
- [Gatheral & Jacquier - Arbitrage-free SVI volatility surfaces](https://arxiv.org/abs/1204.0646)
- [Heston Model Analysis and Implementation](https://arxiv.org/pdf/1502.02963)
- [Full and Fast Calibration of the Heston Model](https://eprints.lse.ac.uk/83754/1/Germano_Full%20and%20fast%20calibration_2017.pdf)
- [Breeden-Litzenberger Risk-Neutral Density Extraction](https://www.newyorkfed.org/medialibrary/media/research/staff_reports/sr677.pdf)
- [FX Volatility Smile Conventions](https://quantpie.co.uk/fx/fx_rr_str.php)

---

### 1. Implied Volatility Extraction

Before calculating skew or smile, we must extract clean implied volatilities from market prices.

#### 1.1 Black-Scholes Implied Volatility

Given market option price `C_mkt`, solve for σ in the Black-Scholes formula:

```
C_BS(S, K, T, r, q, σ) = C_mkt
```

**Newton-Raphson iteration**:
```
σ_{n+1} = σ_n - [C_BS(σ_n) - C_mkt] / Vega(σ_n)
```

Where Vega = ∂C/∂σ = S·e^(-qT)·√T·N'(d₁)

**Implementation details**:
1. Initial guess: σ₀ = √(2π/T) × (C_mkt / S) for ATM options
2. Convergence criterion: |C_BS(σ) - C_mkt| < 1e-8
3. Max iterations: 100 (fail if not converged)
4. Bounds: 0.01 < σ < 5.0 (sanity check)

#### 1.2 Bid-Ask Midpoint vs. Mark

- Use **mid-price** for IV calculation: `(bid + ask) / 2`
- Filter options with spread > 20% of mid (illiquid)
- Weight by inverse spread in calibration: `w_i = 1 / spread_i`

#### 1.3 Moneyness Conventions

We use **log-moneyness** for surface construction:
```
k = ln(K / F)
```
Where F = forward price = S·e^((r-q)T)

**Why log-moneyness?**
- Normalizes across different spot prices
- Natural for stochastic vol models
- Symmetric around ATM (k=0)

---

### 2. Skew Extraction: Delta-Space Analysis

#### 2.1 Black-Scholes Delta Calculation

**Call delta**:
```
Δ_call = e^(-qT) · N(d₁)
```

**Put delta**:
```
Δ_put = -e^(-qT) · N(-d₁)
```

Where:
```
d₁ = [ln(S/K) + (r - q + σ²/2)T] / (σ√T)
```

#### 2.2 Finding Delta-Equivalent Strikes

**Problem**: Market quotes discrete strikes. We need IV at specific deltas (e.g., 25Δ).

**Solution**: Interpolate the IV smile, then solve for strike K such that Δ(K, σ(K)) = target.

**Algorithm**:
```
1. For each expiration T:
   a. Extract all (K_i, IV_i) pairs
   b. Fit interpolator: σ(K) using cubic spline or SVI
   c. For target delta Δ_target (e.g., 0.25 for call, -0.25 for put):
      - Solve: Δ(K, σ(K)) = Δ_target
      - Use root-finding (Brent's method) on K ∈ [K_min, K_max]
   d. Return (K_Δ, σ(K_Δ))
```

**Note on delta conventions**:
- **Spot delta**: Δ = e^(-qT)·N(d₁) — standard for equity options
- **Forward delta**: Δ = N(d₁) — used in some FX markets
- **Premium-adjusted delta**: Accounts for premium paid — used in FX

We use **spot delta** for equity options.

#### 2.3 Risk Reversal Calculation

**25-Delta Risk Reversal (25RR)**:
```
RR_25 = σ_call,25Δ - σ_put,25Δ
```

**Interpretation**:
- RR > 0: Calls are relatively expensive (bullish skew)
- RR < 0: Puts are relatively expensive (bearish skew, typical for equities)

**Standard deltas to track**: 10Δ, 15Δ, 25Δ, 35Δ

#### 2.4 Skew Term Structure

Calculate RR_25 for each expiration, plot against time-to-expiry:

```
Skew term structure: {(T_1, RR_25(T_1)), (T_2, RR_25(T_2)), ...}
```

**Signals**:
- Steep term structure (short-dated RR << long-dated RR): Near-term fear
- Flat term structure: Consistent risk perception across horizons
- Inverted term structure: Unusual, potential mean-reversion opportunity

#### 2.5 Skew Z-Score

Normalize current skew against historical distribution:

```
z_skew = (RR_25,current - μ_RR) / σ_RR
```

Where μ_RR and σ_RR are rolling mean/std (e.g., 252-day window).

**Trade signals**:
- z_skew > 2: Puts extremely rich, consider selling put skew
- z_skew < -2: Puts cheap (rare), consider buying put skew

---

### 3. Smile Extraction: Strike-Space Analysis

#### 3.1 Butterfly Spread as Smile Proxy

**25-Delta Butterfly (25BF)**:
```
BF_25 = [σ_call,25Δ + σ_put,25Δ] / 2 - σ_ATM
```

**Interpretation**:
- BF > 0: Wings are expensive relative to ATM (fat tails priced)
- Higher BF = more convex smile = market expects larger moves

#### 3.2 SVI Parameterization

The **Stochastic Volatility Inspired (SVI)** model parameterizes total implied variance:

```
w(k) = a + b·[ρ·(k - m) + √((k - m)² + σ²)]
```

Where:
- `w(k) = σ²(k)·T` is total variance
- `k = ln(K/F)` is log-moneyness
- Parameters: `{a, b, ρ, m, σ}`

**Parameter interpretation**:
- `a`: Overall variance level
- `b`: Slope of wings (controls how fast IV increases OTM)
- `ρ`: Skew (-1 to 1, negative = put skew)
- `m`: Horizontal translation (shift of smile minimum)
- `σ`: Smoothness at ATM (curvature)

#### 3.3 SVI Calibration Algorithm

**Objective**: Minimize weighted sum of squared errors

```
min Σ w_i · [w_SVI(k_i; θ) - w_mkt(k_i)]²
```

Where w_i = 1/spread_i (liquidity weighting).

**Constraints for arbitrage-free SVI** (Gatheral & Jacquier):
```
a + b·σ·√(1 - ρ²) ≥ 0    (non-negative variance at all strikes)
b ≥ 0                     (wings increase)
|ρ| < 1                   (valid correlation)
σ > 0                     (positive curvature parameter)
```

**Optimization method**: Sequential Least Squares Programming (SLSQP) or L-BFGS-B with bounds.

#### 3.4 Smile Metrics from SVI

Once SVI is calibrated, extract:

1. **ATM volatility**: σ_ATM = √(w(0) / T)
2. **ATM skew**: ∂σ/∂k |_{k=0} = b·ρ / (2·σ_ATM·T)
3. **ATM curvature**: ∂²σ/∂k² |_{k=0} (related to butterfly)
4. **Wing slopes**: lim_{k→±∞} ∂w/∂k = b·(ρ ± 1)

#### 3.5 Smile Asymmetry Index

```
Asymmetry = [σ(k=+0.1) - σ_ATM] - [σ(k=-0.1) - σ_ATM]
```

Positive asymmetry = call wing steeper than put wing (unusual in equities).

---

### 4. Volatility Surface Construction

#### 4.1 SSVI (Surface SVI)

**SSVI** extends SVI to the entire surface by parameterizing ATM variance and skew as functions of time:

```
w(k, T) = θ_T / 2 · {1 + ρ·φ(θ_T)·k + √[(φ(θ_T)·k + ρ)² + (1 - ρ²)]}
```

Where:
- `θ_T = σ²_ATM(T)·T` is ATM total variance at time T
- `φ(θ)` is a function controlling skew decay (e.g., φ(θ) = η / (θ^γ · (1 + θ)^(1-γ)))

**SSVI parameters**: {ρ, η, γ} (global) + {θ_T} for each expiration

#### 4.2 Arbitrage-Free Constraints

**Calendar spread arbitrage**: Total variance must be non-decreasing in T
```
∂w(k, T)/∂T ≥ 0  for all k
```

**Butterfly arbitrage**: Density must be non-negative (Durrleman condition)
```
g(k) = (1 - k·w'/(2w))² - w'/4·(1/w + 1/4) + w''/2 ≥ 0
```

Where w' = ∂w/∂k and w'' = ∂²w/∂k².

#### 4.3 Surface Construction Algorithm

```
Step 1: Data Collection
   - For each expiration T_j, collect (K_i, IV_i, bid_i, ask_i)
   - Filter: remove IV < 5% or IV > 200%, spread > 30%

Step 2: Per-Slice SVI Fit
   - For each T_j, fit SVI to get {a_j, b_j, ρ_j, m_j, σ_j}
   - Check arbitrage constraints; adjust if violated

Step 3: SSVI Global Fit
   - Extract θ_T from SVI fits
   - Fit global {ρ, η, γ} to minimize cross-slice error
   - Iterate: adjust θ_T values to remove calendar arbitrage

Step 4: Interpolation
   - For strikes/expiries not in market data:
   - Use SSVI formula to interpolate
   - Verify butterfly constraint at interpolated points

Step 5: Output
   - Dense grid: strikes from 0.5·S to 2.0·S, expiries from 7d to 2y
   - Store as matrix: IV[strike_idx][expiry_idx]
```

---

### 5. Heston Model Calibration

#### 5.1 Heston Stochastic Differential Equations

**Spot price dynamics**:
```
dS_t = (r - q)·S_t·dt + √v_t·S_t·dW¹_t
```

**Variance dynamics**:
```
dv_t = κ·(θ - v_t)·dt + σ·√v_t·dW²_t
```

**Correlation**:
```
dW¹_t · dW²_t = ρ·dt
```

**Parameters**:
| Parameter | Symbol | Typical Range | Interpretation |
|-----------|--------|---------------|----------------|
| Initial variance | v₀ | 0.01 - 0.25 | Current instantaneous variance |
| Long-term variance | θ | 0.01 - 0.25 | Mean-reversion target |
| Mean-reversion speed | κ | 0.5 - 5.0 | How fast variance reverts |
| Vol of vol | σ | 0.1 - 1.0 | Volatility of variance process |
| Correlation | ρ | -1.0 - 0.0 | Spot-vol correlation (negative = leverage effect) |

**Feller condition** (ensures variance stays positive):
```
2κθ > σ²
```

#### 5.2 Characteristic Function

The Heston model has a semi-closed form solution via the characteristic function:

```
φ(u; T) = exp{C(u, T) + D(u, T)·v₀ + iu·ln(F)}
```

Where F = S·e^((r-q)T) is the forward price, and:

```
C(u, T) = (r - q)·iu·T + (κθ/σ²)·[(κ - ρσiu + d)·T - 2·ln((1 - g·e^(dT))/(1 - g))]

D(u, T) = [(κ - ρσiu + d)/σ²]·[(1 - e^(dT))/(1 - g·e^(dT))]
```

With:
```
d = √[(ρσiu - κ)² + σ²(iu + u²)]
g = (κ - ρσiu + d) / (κ - ρσiu - d)
```

**The "Little Heston Trap"**: The standard formulation can have branch-cut discontinuities. Use the **Albrecher et al. formulation** which selects the correct branch of the complex square root:

```
d = √[(κ - ρσiu)² + σ²(iu + u²)]  (principal branch with Re(d) > 0)
```

#### 5.3 Option Pricing via Fourier Inversion

**Call price** (Carr-Madan / Lewis formulation):

```
C(K, T) = S·e^(-qT)·P₁ - K·e^(-rT)·P₂
```

Where:
```
P_j = 1/2 + (1/π)·∫₀^∞ Re[e^(-iu·ln(K))·φ_j(u)] / (iu) du
```

**Numerical integration**:
- Use **Gauss-Laguerre quadrature** (32-64 points) for semi-infinite integral
- Or **adaptive Simpson's rule** with truncation at u_max ≈ 1000
- **Control variate**: Subtract Black-Scholes integral (known analytically), add back BS price

```python
# Pseudo-code for Heston call price
def heston_call(S, K, T, r, q, v0, theta, kappa, sigma, rho):
    F = S * exp((r - q) * T)
    k = log(K / F)

    # Numerical integration
    integral = gauss_laguerre_integrate(
        lambda u: real(exp(-1j * u * k) * char_func(u, T, params)) / (1j * u),
        n_points=64
    )

    P1 = 0.5 + integral / pi
    P2 = ... # Similar with adjusted char func

    return S * exp(-q * T) * P1 - K * exp(-r * T) * P2
```

#### 5.4 Calibration Objective Function

**Minimize weighted pricing error**:

```
min_θ Σ_i w_i · [C_Heston(K_i, T_i; θ) - C_mkt(K_i, T_i)]²
```

**Alternatively, minimize IV error** (often more stable):
```
min_θ Σ_i w_i · [σ_Heston(K_i, T_i; θ) - σ_mkt(K_i, T_i)]²
```

**Weights**:
- `w_i = vega_i`: More weight on ATM options (higher vega)
- `w_i = 1/spread_i`: More weight on liquid options
- `w_i = 1`: Equal weight (simple)

#### 5.5 Global Optimization

Local optimizers get stuck in local minima. Use **global optimization**:

**Differential Evolution** (recommended):
```
1. Initialize population of N=50 parameter vectors
2. For each generation:
   a. For each vector x_i:
      - Select 3 random vectors x_a, x_b, x_c
      - Create mutant: v = x_a + F·(x_b - x_c), F ∈ [0.5, 1.0]
      - Crossover with x_i to create trial u
      - If f(u) < f(x_i), replace x_i with u
3. Repeat until convergence or max iterations
```

**Parameter bounds**:
```
v0:    [0.001, 1.0]
theta: [0.001, 1.0]
kappa: [0.01, 10.0]
sigma: [0.01, 2.0]
rho:   [-0.99, 0.0]
```

**Regularization** (prevent overfitting):
```
Objective += λ·[(κ - κ_prior)² + (θ - θ_prior)² + ...]
```

Use yesterday's calibrated parameters as priors.

#### 5.6 Calibration Quality Metrics

1. **RMSE** (root mean squared error):
   ```
   RMSE = √[Σ(σ_model - σ_mkt)² / N]
   ```
   Target: RMSE < 0.5% (50 bps)

2. **Max error**:
   ```
   max_i |σ_model(K_i, T_i) - σ_mkt(K_i, T_i)|
   ```
   Target: < 2%

3. **Parameter stability**:
   Track daily parameter changes. Large jumps indicate overfitting or regime change.

---

### 6. Model-Implied Volatility and Relative Value

#### 6.1 Extracting Model-Implied IV

For each option (K, T) in the market:

```
1. Price the option using calibrated Heston: C_Heston(K, T)
2. Invert Black-Scholes to get model-implied IV:
   σ_model such that C_BS(K, T, σ_model) = C_Heston(K, T)
```

This gives us two IV values for each option:
- `σ_mkt`: From market prices
- `σ_model`: From calibrated Heston

#### 6.2 Mispricing Calculation

**Absolute mispricing**:
```
Δσ = σ_mkt - σ_model
```

**Interpretation**:
- Δσ > 0: Market IV > Model IV → Option is **rich** (overpriced)
- Δσ < 0: Market IV < Model IV → Option is **cheap** (underpriced)

**Relative mispricing** (as percentage of model IV):
```
Δσ_rel = (σ_mkt - σ_model) / σ_model × 100%
```

#### 6.3 Statistical Significance

Not all mispricing is tradeable. Test for significance:

**Z-score mispricing**:
```
z_i = Δσ_i / σ_calibration_error
```

Where σ_calibration_error is the standard error from calibration.

**Rule**: Only consider options with |z| > 2 as significantly mispriced.

#### 6.4 Liquidity Adjustment

Adjust mispricing for transaction costs:

```
Δσ_adjusted = Δσ - (spread / 2) / vega
```

Only trade if Δσ_adjusted exceeds threshold (e.g., 1%).

#### 6.5 Relative Value Matrix

Construct a matrix of mispricing by strike and expiration:

```
           T=30d   T=60d   T=90d   T=180d  T=365d
K=0.8·S   +1.2%   +0.8%   +0.5%   +0.3%   +0.1%
K=0.9·S   +0.5%   +0.3%   +0.2%   +0.1%   0.0%
K=S       0.0%    0.0%    0.0%    0.0%    0.0%
K=1.1·S   -0.3%   -0.2%   -0.1%   0.0%    +0.1%
K=1.2·S   -0.8%   -0.5%   -0.3%   -0.1%   +0.2%
```

Positive = rich (sell), Negative = cheap (buy).

---

### 7. Risk-Neutral Density Extraction (Breeden-Litzenberger)

#### 7.1 Theory

The risk-neutral probability density f(K) can be extracted from option prices:

```
f(K) = e^(rT) · ∂²C/∂K²
```

This is model-free (doesn't assume Heston or any specific model).

#### 7.2 Numerical Implementation

**Finite difference approximation**:
```
f(K) ≈ e^(rT) · [C(K+ΔK) - 2·C(K) + C(K-ΔK)] / ΔK²
```

**Steps**:
1. Interpolate call prices across strikes (cubic spline)
2. Compute second derivative numerically
3. Multiply by e^(rT)

**Tail handling**: For extreme strikes where no options trade:
- Fit **Generalized Pareto Distribution** to tails
- Or extrapolate using SVI wing behavior

#### 7.3 Comparing Model vs Market Density

**Heston density**: Inverse Fourier transform of characteristic function:
```
f_Heston(K) = (1/π) · ∫₀^∞ Re[e^(-iu·ln(K)} · φ(u)] du
```

**Comparison**:
- Overlay f_mkt(K) and f_Heston(K)
- Divergences indicate where model fails to capture market expectations
- Large divergence in tails → tail risk mispricing

---

### 8. Temperature-Mapped Relative Value Surface

#### 8.1 Color Scale Definition

Map mispricing Δσ to color using diverging colormap:

```
Color mapping (continuous gradient):
  Δσ ≤ -3%: Deep blue    (#0000FF)  — Very cheap
  Δσ = -2%: Blue         (#4444FF)
  Δσ = -1%: Light blue   (#8888FF)
  Δσ =  0%: White        (#FFFFFF)  — Fair value
  Δσ = +1%: Orange       (#FF8800)
  Δσ = +2%: Red          (#FF4400)
  Δσ ≥ +3%: Deep red     (#FF0000)  — Very rich
```

#### 8.2 Surface Rendering

**3D surface plot**:
- X-axis: Log-moneyness k = ln(K/F) or delta
- Y-axis: Days to expiration (log scale often better)
- Z-axis: Implied volatility
- Color: Relative value (temperature map)

**Implementation** (WebGL/Three.js):
```javascript
// Pseudo-code for surface mesh
for (let i = 0; i < strikes.length; i++) {
    for (let j = 0; j < expiries.length; j++) {
        let iv = surface[i][j].iv;
        let mispricing = surface[i][j].mispricing;

        vertices.push(strikes[i], expiries[j], iv);
        colors.push(mispricingToColor(mispricing));
    }
}
```

#### 8.3 Interactive Features

1. **Hover tooltip**: Show (K, T, σ_mkt, σ_model, Δσ, Greeks)
2. **Slice view**: Fix T, show smile; fix K, show term structure
3. **Time evolution**: Animate surface over historical dates
4. **Filter**: Show only options with |Δσ| > threshold

---

### 9. Trade Identification Algorithm

#### 9.1 Skew Trade: Risk Reversal

**Setup**: Sell rich OTM put, buy cheap OTM call at same delta

**Entry criteria**:
```
1. RR_25 z-score > 2 (puts extremely rich)
2. Put Δσ > +2% (put overpriced)
3. Call Δσ < -1% (call underpriced)
4. Sufficient liquidity (spread < 5%)
```

**Construction**:
- Sell 1x 25Δ put
- Buy 1x 25Δ call
- Net delta ≈ +50Δ (bullish bias)

**Exit**:
- RR mean-reverts (z-score < 0.5)
- Or stop-loss at RR z-score > 3.5

#### 9.2 Ratio Spread (Gamma Picking)

**Setup**: Sell 1 rich put, buy N cheap calls (N = 3-4)

**Entry criteria**:
```
1. Put Δσ > +3% (significantly rich)
2. Call Δσ < -2% (significantly cheap)
3. Funding: put premium ≥ N × call premium
4. Max loss acceptable at put strike
```

**Construction** (example: 1x4 ratio):
- Sell 1x 20Δ put @ $5.00 premium
- Buy 4x 10Δ calls @ $1.25 each
- Net credit: $5.00 - $5.00 = $0 (or small credit)

**P&L profile**:
- Below put strike: Max loss = (put strike - spot) - net credit
- Between strikes: Profit = net credit + time decay
- Above call strike: Unlimited upside (convex payoff)

#### 9.3 Butterfly (Smile Trade)

**Setup**: Exploit smile curvature mispricing

**Entry criteria**:
```
1. BF_25 z-score > 2 (wings extremely rich)
2. Wing options Δσ > +2%
3. ATM options fairly valued (|Δσ| < 0.5%)
```

**Construction**:
- Sell 1x 25Δ call + Sell 1x 25Δ put (short wings)
- Buy 2x ATM straddle (hedge)
- Net: Short butterfly (profit if realized vol < implied curvature)

---

### 10. Daily Workflow

```
08:00 - Market Open Prep
├── Fetch overnight index futures, VIX levels
├── Load previous day's calibrated parameters
└── Check for corporate actions, earnings dates

09:30 - Market Open
├── Begin streaming options chain data
├── Calculate real-time IVs as quotes update
└── Monitor ATM IV vs previous close

10:00 - First Calibration
├── Run Heston calibration on first 30min of data
├── Generate initial relative value surface
├── Flag any options with |Δσ| > 2%
└── Send alerts for potential trades

Hourly - Rolling Updates
├── Recalibrate Heston every hour (or on significant moves)
├── Update relative value matrix
├── Refresh 3D visualization
└── Log parameter drift

15:30 - End of Day
├── Final calibration snapshot
├── Store all surfaces, parameters, signals
├── Generate daily report:
│   ├── Calibration quality (RMSE, max error)
│   ├── Parameter changes from prior day
│   ├── Top 10 mispriced options (rich and cheap)
│   └── Trade signals triggered
└── Archive for backtesting

Overnight
├── Run backtests on historical signals
├── Compute strategy performance metrics
└── Update momentum signals for next day
```

---

### 11. Backtesting Framework

#### 11.1 Signal Generation (Historical)

For each historical date t:
1. Load options chain snapshot
2. Calibrate Heston to that day's surface
3. Calculate relative value for all options
4. Generate trade signals per rules above

#### 11.2 Trade Simulation

```
For each signal:
1. Entry:
   - Record entry price (mid or aggressive)
   - Calculate position Greeks

2. Daily mark-to-market:
   - Reprice position using next day's surface
   - Track P&L, Greeks evolution

3. Exit:
   - Exit when signal reverses, or at expiration
   - Record final P&L, holding period
```

#### 11.3 Performance Metrics

| Metric | Formula | Target |
|--------|---------|--------|
| Win rate | # winners / # trades | > 55% |
| Avg win / Avg loss | Mean(winners) / Mean(losers) | > 1.5 |
| Sharpe ratio | Mean(returns) / Std(returns) × √252 | > 1.5 |
| Max drawdown | Max peak-to-trough decline | < 15% |
| Calmar ratio | CAGR / Max drawdown | > 1.0 |

---

### 12. Risk Management

#### 12.1 Position Limits

- **Single name**: Max 5% of capital per underlying
- **Delta**: Net delta < 20% of portfolio
- **Vega**: Net vega < 2% of portfolio (vol exposure)
- **Gamma**: Monitor gamma; reduce before large events

#### 12.2 Scenario Analysis

For each position, calculate P&L under:
- Spot ±5%, ±10%, ±20%
- IV ±5 vol points, ±10 vol points
- Combined: spot down 10%, IV up 10 points (crash scenario)

#### 12.3 Greeks Aggregation

Real-time portfolio Greeks:
```
Portfolio Δ = Σ position_i × Δ_i
Portfolio Γ = Σ position_i × Γ_i
Portfolio V = Σ position_i × V_i (vega)
Portfolio Θ = Σ position_i × Θ_i
```

Alert if any Greek exceeds threshold.
