package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Hash     []byte // hash of the current block
	Data     []byte // transaction data
	PrevHash []byte // hash of previous block
}

type BlockChain struct {
	blocks []*Block
}

func (b *Block) DeriveHash() {
	hash := sha256.Sum256(bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})) // fixed output size of 32 bytes regardless of input
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

func initBlockChain() *BlockChain {
	return &BlockChain{[]*Block{CreateBlock("Genesis", []byte{})}}
}
func main() {
	bchain := initBlockChain()

	bchain.AddBlock("First Block")
	bchain.AddBlock("Second Block")
	bchain.AddBlock("Third Block")

	for _, block := range bchain.blocks {
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
