---
name: blockchain-contract-attack
description: >-
  Blockchain/smart contracts: Etherscan, slither/mythril, reentrancy, access control, oracles, flash loans, cross-chain bridges, and exposed RPC endpoints. Use when auditing smart contracts, DeFi, or blockchain attack surfaces.
metadata:
  tags: [penetration-testing, red-team]
---

## Blockchain / Smart Contracts

```
=== Blockchain / smart contracts ===
Source: Etherscan getsourcecode API | Audit with slither/mythril/manticore
Manual review: reentrancy (.call{value}: transfer before state update violates C-E-I) | access control (missing onlyOwner) | integer overflow (<0.8 without SafeMath)
  oracle manipulation (instant AMM-price manipulation with a flash loan) | randomness (block.timestamp is controllable) | authorization abuse (unlimited approve/permit replay) | delegatecall proxy storage collisions
DeFi: flash-loan attacks / sandwich attacks (front-running) / governance attacks / signature replay | cross-chain bridges: signature-threshold bypass / replay | RPC: exposed 8545 allowing direct eth_sendTransaction
```
