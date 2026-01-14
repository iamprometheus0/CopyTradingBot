+----------------------+----------------------+----------------------+
| \                    | Recommendation       | Rationale & Notes    |
| \                    |                      |                      |
| \                    |                      |                      |
| \                    |                      |                      |
| Below is the         |                      |                      |
| architecture for my  |                      |                      |
| Polymarket           |                      |                      |
| copy-trading bot.    |                      |                      |
| This design targets  |                      |                      |
| copying              |                      |                      |
| high-frequency bots  |                      |                      |
| (multiple            |                      |                      |
| orders/minute) in    |                      |                      |
| volatile 15-minute   |                      |                      |
| crypto Up/Down       |                      |                      |
| binary markets (BTC, |                      |                      |
| ETH, SOL, XRP), with |                      |                      |
| strong emphasis on   |                      |                      |
| minimizing           |                      |                      |
| end-to-end latency   |                      |                      |
| to reduce slippage:\ |                      |                      |
| \                    |                      |                      |
| High-Level           |                      |                      |
| Architecture         |                      |                      |
| Overview\            |                      |                      |
| The bot follows a    |                      |                      |
| clean, modular       |                      |                      |
| event-driven         |                      |                      |
| design:\             |                      |                      |
| \                    |                      |                      |
| Real-time Monitoring |                      |                      |
| → Detect target      |                      |                      |
| bot/trader activity  |                      |                      |
| instantly.\          |                      |                      |
| Decision & Filtering |                      |                      |
| → Validate (market   |                      |                      |
| type, size           |                      |                      |
| threshold, slippage  |                      |                      |
| check, etc.).\       |                      |                      |
| Order Execution →    |                      |                      |
| Sign & submit        |                      |                      |
| proportional copy    |                      |                      |
| orders fast.\        |                      |                      |
| State & Safety →     |                      |                      |
| Track positions,     |                      |                      |
| balances, PnL;       |                      |                      |
| enforce risk rules.\ |                      |                      |
| Logging & Monitoring |                      |                      |
| → Latency            |                      |                      |
| histograms, alerts   |                      |                      |
| on high slippage.\   |                      |                      |
| \                    |                      |                      |
| Key principles:\     |                      |                      |
| \                    |                      |                      |
| Push-based           |                      |                      |
| (WebSockets) over    |                      |                      |
| polling.\            |                      |                      |
| Low-latency infra    |                      |                      |
| (VPS close to        |                      |                      |
| Polymarket           |                      |                      |
| routing).\           |                      |                      |
| Concurrency via      |                      |                      |
| goroutines for       |                      |                      |
| parallel WS          |                      |                      |
| handling +           |                      |                      |
| execution.\          |                      |                      |
| Non-custodial ---    |                      |                      |
| your private key     |                      |                      |
| signs orders         |                      |                      |
| locally.\            |                      |                      |
| \                    |                      |                      |
| Recommended Software |                      |                      |
| Stack\               |                      |                      |
| \                    |                      |                      |
| **Component**        |                      |                      |
+----------------------+----------------------+----------------------+
| Language             | **Go 1.23+**         | Fast compilation,    |
|                      |                      | excellent            |
|                      |                      | goroutines/channels  |
|                      |                      | for concurrency, low |
|                      |                      | GC impact. Quick     |
|                      |                      | ramp-up from         |
|                      |                      | TS/Node.js.          |
+----------------------+----------------------+----------------------+
| **WebSocket          | \- **CLOB WS**:      | Official endpoints   |
| Clients**            | wss:/                | for sub-100 ms       |
|                      | /ws-subscriptions-cl | pushes. Use          |
|                      | ob.polymarket.com/ws | community Go clients |
|                      | (market + user       | (e.g.,               |
|                      | channels) -          | github.com/Matthew1  |
|                      | **RTDS**:            | 7-21/go-polymarket-r |
|                      | wss://ws-live        | eal-time-data-client |
|                      | -data.polymarket.com | or                   |
|                      | (crypto prices,      | github.com/ivanzz    |
|                      | comments)            | eth/polymarket-go-re |
|                      |                      | al-time-data-client) |
|                      |                      | for RTDS. For CLOB,  |
|                      |                      | implement with       |
|                      |                      | gorilla/websocket or |
|                      |                      | n                    |
|                      |                      | hooyr.io/websocket + |
|                      |                      | auto-reconnect.      |
+----------------------+----------------------+----------------------+
| **Order Signing &    | Official             | Handles EIP-712      |
| CLOB**               | **go-order-utils**   | signing for          |
|                      | (from Polymarket     | limit/market orders. |
|                      | docs) or community   | Use proxy wallet     |
|                      | like                 | setup for safety.    |
|                      | gi                   |                      |
|                      | thub.com/mtt-labs/po |                      |
|                      | ly-market-sdk/client |                      |
+----------------------+----------------------+----------------------+
| **Polygon / Wallet** | **go-ethereum**      | For any on-chain     |
|                      | (github.com/e        | reads (balances,     |
|                      | thereum/go-ethereum) | allowances), but     |
|                      | or **viem-go**       | minimize --- most    |
|                      | equivalent           | action is off-chain  |
|                      |                      | CLOB.                |
+----------------------+----------------------+----------------------+
| **Data / State**     | **Redis**            | Redis: fast caching  |
|                      | (in-memory,          | of positions, recent |
|                      | github.c             | trades, rate limits. |
|                      | om/redis/go-redis) + | Postgres: persistent |
|                      | optional             | history/PnL logs.    |
|                      | **PostgreSQL**       |                      |
+----------------------+----------------------+----------------------+
| **Config & Env**     | **viper** or         | Easy YAML/JSON + env |
|                      | **koanf** + .env     | vars for targets,    |
|                      |                      | multipliers,         |
|                      |                      | filters.             |
+----------------------+----------------------+----------------------+
| **Logging**          | **zerolog** or       | Structured, fast,    |
|                      | **zap**              | with latency fields. |
+----------------------+----------------------+----------------------+
| **Metrics /          | **promethe           | Track                |
| Monitoring**         | us/client_golang** + | detection-to-fill    |
|                      | **Grafana**          | latency, slippage %, |
|                      | (optional)           | win rate. Alert via  |
|                      |                      | Telegram if    |
|                      |                      | latency \>500 ms.    |
+----------------------+----------------------+----------------------+
| **Deployment**       | **Docker** + single  | Easy to build, run,  |
|                      | container (or        | restart. Use         |
|                      | systemd on VPS)      | multi-stage          |
|                      |                      | Dockerfile for small |
|                      |                      | image.               |
+----------------------+----------------------+----------------------+
| **Infra / VPS**      | **NYC / New Jersey** | East Coast US edges  |
|                      | location preferred   | out for lowest ping  |
|                      | (e.g.,               | to                   |
|                      | NinjaMobileTrader,   | P                    |
|                      | QuantVPS NYC/NJ,     | olymarket/Cloudflare |
|                      | Vultr HF NY/NJ) →    | routing. Test ping   |
|                      | **1--5 ms** ping     | to                   |
|                      | claims.              | clob.polymarket.com  |
|                      | A                    | & ws endpoints.      |
|                      | msterdam/Netherlands | Start with 4--8 GB   |
|                      | fallback             | RAM, 2--4 vCPU, SSD  |
|                      | (TradingVPS.io,      | (\~\$15--40/mo).     |
|                      | QuantVPS Amsterdam)  | Avoid Mumbai/home    |
|                      | → **10--30 ms**,     | (150+ ms).           |
|                      | often                |                      |
|                      | Cloudflare-friendly. |                      |
+----------------------+----------------------+----------------------+

\
\
\
**Component Breakdown & Data Flow**

  **Startup / Config**

    -   Load env: private key (hot wallet), target wallets (array),
        trade multiplier (e.g., 0.3×), filters (min size, max slippage
        %, only 15-min crypto markets via token IDs or question regex).

    -   Connect to Redis/Postgres.

```{=html}
<!-- -->
```
  **WebSocket Listeners** (goroutines)

    -   **CLOB WS** → Subscribe to /market channel (public
        orderbook/prices/trades). Filter for activity on 15-min Up/Down
        markets.

    -   **RTDS WS** → Subscribe to crypto_prices (e.g., BTCUSDT) for
        spot reference (helps detect mispricings or confirm momentum).

    -   **Target wallet monitoring** → If CLOB WS doesn\'t directly push
        per-user trades (it does for your own user channel), use Data
        API fallback polling sparingly (\~every 500--1000 ms) or watch
        public trades and match wallet addresses.

    -   On message → parse JSON → if new trade from target wallet in
        relevant market → push to decision channel.

```{=html}
<!-- -->
```
1.  **Decision Engine** (goroutine consuming channel)

    # Polymarket Copy-Trading Bot – Refined Decision Engine Logic v2 (With 3% Fill Price Deviation Cap)

Goal: Faithful mirroring of high-frequency bots in 15-minute crypto Up/Down binary markets, while protecting against unacceptable slippage.  
Enforce a strict **3% adverse deviation cap** between detected target fill price and your prospective fill price.

## Core Rules
- Mirror trade-for-trade as closely as possible.
- **Skip** the copy if current best available price deviates **>3% adversely** from the target bot's detected fill price.
- If within ≤3% deviation → attempt fill at the **next best price** (aggressive limit or market order to cross spread if needed).
- Do **not** force maker-only orders — prioritize fill speed and closeness to target price.
- Preserve mechanical repetition (accumulation, hedging sequences) while avoiding catastrophic slippage.

## Decision Engine Flow

1. WS Event Received (CLOB trade push or RTDS activity)
   - Filter: 15-min crypto binary market? From target wallet?
   - If no → discard

2. Extract Target Fill Price
   - detectedPrice = event.fillPrice (or last matched price from trade event)

3. Fetch Current Order Book (real-time via CLOB WS or quick GET)
   - For buy copy: bestAsk = lowest ask price
   - For sell copy: bestBid = highest bid price
   - prospectivePrice = bestAsk (if buying) or bestBid (if selling)

4. Deviation Check (adverse only)
   - deviationPct = ((prospectivePrice - detectedPrice) / detectedPrice) * 100   (for buys)
   - deviationPct = ((detectedPrice - prospectivePrice) / detectedPrice) * 100   (for sells)
   - If deviationPct > +3.0% (worse than target) → log "Excessive slippage" + skip this trade
   - If deviationPct ≤ +3.0% → proceed

5. Market Stability Guard (Rare-Event Fuse)
    Empirical observation shows that 15-min and 1-hour crypto binary markets on Polymarket typically maintain ~$0.01 spreads under normal conditions. A sudden spread expansion (>3% of mid) therefore indicates an abnormal or black-swan regime (liquidity pull, sweep, or repricing). In such cases, copy execution is skipped to avoid interacting with post-impact, discontinuous liquidity.

6. Build Order at Next Best Price
   - If prospectivePrice == detectedPrice → use market/taker order (fastest fill)
   - If prospectivePrice ≠ detectedPrice but within 3% → place aggressive limit order:
     - Buy: limitPrice = bestAsk + small buffer (e.g., +0.1–0.3% to cross spread)
     - Sell: limitPrice = bestBid - small buffer
     - Set short TTL (e.g., 5–10 seconds) to avoid stale orders
   - Fallback: If limit unlikely to fill quickly → use market order (accept crossing spread)

7. Minimal Safety Checks (before submit)
   - Sufficient USDC for size × multiplier?
   - Per-market exposure cap not exceeded (e.g., max 10–20% capital per slot)?
   - Trade size > noise threshold ($50–100 equivalent)?
   - If any fail → log & skip

8. Submit Async
   - Sign with go-order-utils (EIP-712)
   - POST to clob.polymarket.com/orders
   - Run in goroutine pool

9. Post-Execution Monitoring
   - Record actual fill price & size in Redis
   - Calculate realized slippage % vs. detectedPrice
   - Alert (Telegram/Slack) if actual fill >3% worse than detected (post-facto visibility)
   - Track blended avg cost per side (YES/NO) per market
   - Alert if hedge becomes unbalanced (Qty_YES / Qty_NO > 2× deviation)
```{=html}
<!-- -->
```
  **Execution** (goroutine pool or dedicated)

    -   Sign order with go-order-utils (EIP-712).

    -   POST to
        [[https://clob.polymarket.com/orders](https://clob.polymarket.com/orders?referrer=grok.com)]{.underline}.

    -   On success → record in Redis (position tracking).

    -   Retry on transient fails (exponential backoff).

```{=html}
<!-- -->
```
  **Position & Risk Management** (separate goroutine)

    -   Periodically sync your positions via Gamma/Data API or CLOB user
        WS.

    -   Auto-exit logic if needed (e.g., max drawdown, time-based on
        15-min resolve).

    -   Circuit breaker: pause copying if latency spikes or vol too
        high.

```{=html}
<!-- -->
```
  **Latency Critical Path**

    -   Goal: **Detection** \<100 ms (WS push) + **Decision** \<10 ms +
        **Sign/Submit** \<100 ms + **Fill** \<200 ms → total
        **\~200--500 ms**.

    -   Measure & log every step → tune VPS if \> expected.
