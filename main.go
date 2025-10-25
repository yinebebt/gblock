package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"time"
)

// Block represents a single block in the blockchain.
type Block struct {
	Hash      []byte
	Data      []byte
	PrevHash  []byte
	Timestamp int64
	Nonce     int
}

// BlockChain represents a chain of blocks.
type BlockChain struct {
	blocks     []*Block
	Difficulty int
}

// DeriveHash calculates and sets the block's hash.
func (b *Block) DeriveHash() {
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(b.Timestamp))

	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(b.Nonce))

	info := bytes.Join([][]byte{b.Data, b.PrevHash, timestampBytes, nonceBytes}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// Mine performs Proof-of-Work by finding a valid nonce.
func (b *Block) Mine(difficulty int) {
	target := strings.Repeat("0", difficulty)
	var hashStr string
	attempts := 0

	for {
		b.DeriveHash()
		hashStr = fmt.Sprintf("%x", b.Hash)

		if hashStr[:difficulty] == target {
			fmt.Printf("Block mined. Nonce: %d, Attempts: %d\n", b.Nonce, attempts)
			fmt.Printf(" Hash: %s\n", hashStr)
			break
		}

		b.Nonce++
		attempts++

		if attempts%10_000 == 0 {
			fmt.Printf("Mining... attempts: %d, current hash: %s\n", attempts, hashStr[:16])
		}
	}
}

func createBlock(data string, prevHash []byte, difficulty int) *Block {
	block := &Block{
		Hash:      []byte{},
		Data:      []byte(data),
		PrevHash:  prevHash,
		Timestamp: time.Now().Unix(),
		Nonce:     0,
	}

	fmt.Printf("Mining block: %s\n", data)
	block.Mine(difficulty)

	return block
}

// ValidateBlock verifies the block's hash and difficulty.
func (b *Block) ValidateBlock(difficulty int) bool {
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(b.Timestamp))

	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(b.Nonce))

	info := bytes.Join([][]byte{b.Data, b.PrevHash, timestampBytes, nonceBytes}, []byte{})
	recalculatedHash := sha256.Sum256(info)

	if !bytes.Equal(b.Hash, recalculatedHash[:]) {
		return false
	}

	hashStr := fmt.Sprintf("%x", b.Hash)
	target := strings.Repeat("0", difficulty)

	return hashStr[:difficulty] == target
}

// ValidateChain verifies the integrity of the entire blockchain.
func (chain *BlockChain) ValidateChain() bool {
	for i := 1; i < len(chain.blocks); i++ {
		currentBlock := chain.blocks[i]
		prevBlock := chain.blocks[i-1]

		if !currentBlock.ValidateBlock(chain.Difficulty) {
			log.Printf("Block %d has invalid hash or doesn't meet difficulty", i)
			return false
		}

		if !bytes.Equal(currentBlock.PrevHash, prevBlock.Hash) {
			log.Printf("Block %d is not properly linked to block %d", i, i-1)
			return false
		}

		if currentBlock.Timestamp < prevBlock.Timestamp {
			log.Printf("Block %d has invalid timestamp", i)
			return false
		}
	}

	if !chain.blocks[0].ValidateBlock(chain.Difficulty) {
		log.Println("Genesis block has invalid hash or doesn't meet difficulty")
		return false
	}

	return true
}

// AddBlock creates and appends a new block to the chain.
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := createBlock(data, prevBlock.Hash, chain.Difficulty)
	chain.blocks = append(chain.blocks, newBlock)
}

func initBlockChain(difficulty int) *BlockChain {
	return &BlockChain{
		blocks:     []*Block{createBlock("Genesis", []byte{}, difficulty)},
		Difficulty: difficulty,
	}
}

func main() {
	difficulty := 3

	fmt.Println("Proof-of-Work Blockchain")
	fmt.Printf("\nDifficulty: %d\n\n", difficulty)

	bchain := initBlockChain(difficulty)

	fmt.Println("\nMining Block 1")
	bchain.AddBlock("Alice sends 10 BTC to Bob")

	fmt.Println("\nMining Block 2")
	bchain.AddBlock("Bob sends 5 BTC to Charlie")

	fmt.Println("\nMining Block 3")
	bchain.AddBlock("Charlie sends 3 BTC to Alice")

	fmt.Println("\nBlockchain Contents")
	for i, block := range bchain.blocks {
		fmt.Printf("\n--- Block %d ---\n", i)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
	}

	fmt.Println("\nValidation")
	if bchain.ValidateChain() {
		fmt.Println("Chain is valid")
	} else {
		fmt.Println("Chain is invalid")
	}

	fmt.Println("\nTamper Detection")
	oldHash := fmt.Sprintf("%x", bchain.blocks[1].Hash)
	bchain.blocks[1].Data = []byte("Alice sends 1000 BTC to Bob")

	if bchain.ValidateChain() {
		fmt.Println("Chain is valid")
	} else {
		fmt.Println("Chain is invalid - tampering detected")
	}

	fmt.Printf("\nStored hash:   %s\n", oldHash)

	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(bchain.blocks[1].Timestamp))
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(bchain.blocks[1].Nonce))
	info := bytes.Join([][]byte{bchain.blocks[1].Data, bchain.blocks[1].PrevHash, timestampBytes, nonceBytes}, []byte{})
	recalculatedHash := sha256.Sum256(info)

	fmt.Printf("Actual hash:   %x\n", recalculatedHash)

	fmt.Println("\nDifficulty Comparison")
	for d := 1; d <= 6; d++ {
		avgAttempts := 1 << (4 * uint(d))
		fmt.Printf("Difficulty %d: ~%d attempts\n", d, avgAttempts)
	}
}
