# Blockchain Concepts

Core concepts implemented in this project.

## Cryptographic Hashing

A cryptographic hash function is a one-way mathematical function with these properties:

1. Deterministic - same input always produces same output
2. Fixed-size output - SHA-256 produces 256 bits
3. Irreversible - cannot derive input from output
4. Avalanche effect - small input change produces completely different output
5. Collision resistant - computationally infeasible to find two inputs with same output

### SHA-256 Example

```
Input: "Hello"
Hash: 185f8db32271fe25f561a6fc938b2e264306ec304eda518007d1764826381969

Input: "Hello!"
Hash: 334d016f755cd6dc58c53a86e183882f8ec14f52fb05345887c8a5edd42c87b7
```

One character change produces completely different hash.

## Block Structure

```go
type Block struct {
    Hash      []byte
    Data      []byte
    PrevHash  []byte
    Timestamp int64
    Nonce     int
}
```

**Hash**: SHA-256(Data + PrevHash + Timestamp + Nonce)
**Data**: Transaction or message content
**PrevHash**: Links to previous block (creates chain)
**Timestamp**: Unix timestamp in seconds
**Nonce**: Number used once (for Proof-of-Work)

## Chain Linking

Each block contains the hash of the previous block, creating an immutable chain:

```
Block 0 (Genesis)     Block 1              Block 2
Hash: 0xABCD...  -->  Hash: 0x1234...  --> Hash: 0x5678...
PrevHash: 0x0000      PrevHash: 0xABCD     PrevHash: 0x1234
```

### Tamper Evidence

If Block 1 data is modified:

- Block 1 hash changes
- Block 2 PrevHash no longer matches Block 1 hash
- Chain breaks

Attacker must re-mine all subsequent blocks, which is computationally expensive.

## Timestamps

Unix timestamps provide:

- Event ordering
- Audit trail
- Consensus on time-ordering
- Protection against timestamp manipulation

Validation checks:

- Block timestamp >= previous block timestamp
- Block timestamp <= current time + allowed drift

## Validation

### Block Validation

1. Recalculate hash from block contents
2. Verify recalculated hash matches stored hash
3. Verify hash meets difficulty requirement (leading zeros)

### Chain Validation

1. Validate each block's hash
2. Verify each block links to previous block (PrevHash == prevBlock.Hash)
3. Verify timestamp ordering
4. Validate genesis block separately

## Proof-of-Work (PoW)

Consensus mechanism that makes block creation computationally expensive but verification cheap.

### The Double-Spend Problem

Without consensus, a distributed system has no way to agree on which transactions are valid:

```
Alice has 10 BTC
Alice sends 10 BTC to Bob     (Block A)
Alice sends 10 BTC to Charlie (Block B)
```

Which is valid? Consensus solves this.

### Mining Process

1. Collect transactions
2. Create block with Data, PrevHash, Timestamp, Nonce=0
3. Calculate hash
4. If hash starts with N zeros (meets difficulty), block is valid
5. Otherwise, increment nonce and try again
6. Broadcast valid block to network

### Nonce

A counter that miners increment to find a hash meeting the difficulty target:

```
Nonce: 0     Hash: 8f3a2b... (invalid)
Nonce: 1     Hash: 7c1d9e... (invalid)
Nonce: 2     Hash: 6b8f2a... (invalid)
...
Nonce: 47853 Hash: 0000a1b2... (valid - starts with 4 zeros)
```

### Difficulty

Number of leading zeros required in hash. Each additional zero makes mining ~16x harder:

```
Difficulty 1: ~16 attempts
Difficulty 2: ~256 attempts
Difficulty 3: ~4,096 attempts
Difficulty 4: ~65,536 attempts
```

Bitcoin uses ~19 leading zeros = ~10 minutes per block.

### Mining Algorithm

```go
func (b *Block) Mine(difficulty int) {
    target := strings.Repeat("0", difficulty)

    for {
        b.DeriveHash()
        hashStr := fmt.Sprintf("%x", b.Hash)

        if hashStr[:difficulty] == target {
            break
        }

        b.Nonce++
    }
}
```

## PoW Properties

**Security:**

- Immutability - changing old blocks requires re-mining all subsequent blocks
- Sybil resistance - cannot fake computational work
- Fair leader election - anyone with hash power can participate
- Objective consensus - longest chain (most cumulative work) wins

**Tradeoffs:**

- High energy consumption
- Slow finality (need multiple confirmations)
- Vulnerable to 51% attack if hash power centralizes
- Requires expensive mining hardware

## Consensus Comparison

| Mechanism | Security Model      | Energy | Speed     | Examples                      |
|-----------|---------------------|--------|-----------|-------------------------------|
| PoW       | Computational work  | High   | Slow      | Bitcoin, Ethereum (pre-merge) |
| PoS       | Economic stake      | Low    | Fast      | Ethereum, Cardano             |
| PoA       | Trusted validators  | Low    | Very fast | Private chains                |
| BFT       | Byzantine agreement | Low    | Fast      | Tendermint, Cosmos            |
