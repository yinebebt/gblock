# gblock

A minimal blockchain implementation in Go.

### Core Components
- **Block**:  
  - `Data` (content), 
  - `PrevHash` (previous block's SHA-256), 
  - `Hash` (SHA-256 of `Data+PrevHash`).
- **Blockchain**: Slice of linked `Block` structs.

### Key Logic
- `DeriveHash()`: Generates SHA-256 hash from `Data + PrevHash`.
- `CreateBlock()`: Initializes a block and computes its hash.
- Genesis block starts the chain with `PrevHash = []byte{}`.
