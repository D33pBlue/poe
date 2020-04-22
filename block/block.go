/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: block.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



// Package block contains the definitions
// of blocks, transactions, nonce
package block

// A block is a collection of transactions,
// with a nonce and the hash-link of the previous one.
type Block struct{
  Transactions []Transact
//  Nonce Nonce
  HashPrev []byte
}

// Store a transaction inside the block, after checking
// its validity through Transact.Check()
func (self *Block)AddTransaction(t Transact){
  if t.Check() {
    self.Transactions = append(self.Transactions,t)
  }
}

// Returns the current hash of the block (without
// the nonce if withoutNonce==true, complete otherwise).
func (self *Block)GetHash(withoutNonce bool)[]byte{
  return nil
}
