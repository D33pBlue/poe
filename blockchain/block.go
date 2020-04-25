/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: block.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-25
 * @Copyright: 2020
 */

// Package block contains the definitions
// of blocks, transactions, nonce
package blockchain

import(
  "time"
  // "github.com/D33pBlue/poe/utils"
)

type Block struct{
  Previous *Block
  LenSubChain int
  Transactions *Tree
  Timestamp time.Time
  NumJobs int
  NonceNoJob NonceNoJob
  MiniBlocks []MiniBlock
  Hash []byte
}

// The Check method checks the validity of the block.
// With deep==true this method performs the check in details,
// considering all the transactions.
// With deep==false, instead, the miner can
// skip the checks of his own data.
func (self *Block)Check(deep bool)bool{
  return true // TODO: implement later
}

func (self *Block)GetTransaction(hash []byte)Transaction{
  return nil // TODO: implement later
}
