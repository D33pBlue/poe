/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: block.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-23
 * @Copyright: 2020
 */

// Package block contains the definitions
// of blocks, transactions, nonce
package block

import(
  "github.com/D33pBlue/poe/transact"
)

type Block struct{
  // TODO: implement later
}

func (self *Block)GetTransaction(hash []byte)transact.Transaction{
  return nil // TODO: implement later
}
