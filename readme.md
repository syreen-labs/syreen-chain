# Syreen Chain

A high-performance Layer 1 blockchain built for decentralized trading. On-chain DEX, ~375ms blocks, 80% fee burn, EVM compatible, IBC cross-chain.

- **Website:** https://syreen.com
- **Explorer:** https://explorer.syreen.com
- **Wallet:** https://wallet.syreen.com
- **Whitepaper:** https://syreen.com/syreen-whitepaper.pdf

## Key Features

- **On-Chain DEX** — AMM + order book with limit, stop-loss, and take-profit orders
- **~375ms Block Time** — Sub-second finality via CometBFT consensus
- **80% Fee Burn** — Deflationary tokenomics, every transaction reduces supply
- **EVM Compatible** — Deploy Solidity contracts (Chain ID: 79733)
- **IBC Cross-Chain** — Connected to the Cosmos ecosystem
- **Intent-Based Trading** — Solver network finds optimal execution paths
- **MEV Protection** — Protocol-level commit-reveal ordering
- **Copy Trading** — Follow successful traders with on-chain track records

## Tokenomics

| Allocation | % | Amount |
|------------|---|--------|
| Liquidity Pools | 30% | 30M SYR |
| Trading Rewards | 20% | 20M SYR |
| Founder / Team | 15% | 15M SYR (vested) |
| Treasury | 10% | 10M SYR |
| Marketing / Listings | 10% | 10M SYR |
| Community / Airdrops | 9.6% | 9.6M SYR |
| Staking Reserve | 5% | 5M SYR |
| Genesis Validators | 0.4% | 400K SYR |
| **Total** | **100%** | **100M SYR** |

Founder tokens are subject to a 6-month cliff + 24-month linear vest, enforced on-chain via `PeriodicVestingAccount`.

## Build from Source

Requires Go 1.22+.

```bash
# Clone
git clone https://github.com/syreen-labs/syreen-chain.git
cd syreen-chain

# Build (static binary)
CGO_ENABLED=1 go build -buildvcs=false -tags "netgo" \
  -ldflags '-linkmode external -extldflags "-static"' \
  -trimpath -o syreend ./cmd/syreend
```

## Network Info

| Parameter | Value |
|-----------|-------|
| Chain ID (Cosmos) | syreen-1 |
| Chain ID (EVM) | 79733 |
| Bond Denom | usyreen (1 SYR = 1,000,000 usyreen) |
| Block Time | ~375ms |
| Max Validators | 150 |
| Inflation | 5-15% adaptive |
| Fee Burn | 80% |
| Genesis Supply | 100,000,000 SYR |

## Built With

- [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) v0.53.5
- [CometBFT](https://github.com/cometbft/cometbft) v0.38.21
- [IBC-Go](https://github.com/cosmos/ibc-go) v10.4.0

## Community

- [Discord](https://discord.gg/HH5vexkEq)
- [Twitter / X](https://x.com/SyreenChain)
- [Telegram](https://t.me/syreenchain)

## License

Copyright 2026 Syreen Labs. All rights reserved.
