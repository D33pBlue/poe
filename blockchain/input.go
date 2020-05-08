/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: input.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-08
 * @Copyright: 2020
 */


package blockchain

import(
  "strconv"
)

// A TrInput represents inside a transaction the source of the
// money a client is spending. It points to a TrOutput that is
// included in a transaction, that is stored inside a block
// in the blockchain.
type TrInput struct{
  Block string // hash of the old block with the transaction to spend
  ToSpend string // hash of the old transaction to spend
  Index int // index inside ToSpend of the record to spend
  // proof of ownership?
}


// Returns the TrOutput this TrInput points to, if it
// exists in the chain pointed by the head block; nil otherwise.
func (self *TrInput)GetSource(head *Block)*TrOutput{
  block := head.FindPrevBlock(self.Block)
  if block==nil{ return nil }
  transaction := block.FindTransaction(self.ToSpend)
  if transaction==nil{ return nil }
  return transaction.GetOutputAt(self.Index)
}

func (self *TrInput)ToString()string{
  return self.Block+","+self.ToSpend+","+strconv.Itoa(self.Index)
}
