/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: input.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-07
 * @Copyright: 2020
 */


package blockchain

import(
  "github.com/D33pBlue/poe/utils"
)

// A TrOutput represents inside a transaction the receiver of
// some money. This output can be spent in another transaction.
// The hash of the block in which this output is used inside
// a transaction (it is spent) can be stored or retrieved
// through SetSpentIn and GetSpentIn methods.
type TrOutput struct{
  Address utils.Addr // address of the receiver
  Value int // value to exchange
  spent string
}

// Returns the hash of the block in which this
// output has been spent, or empty string if it
// is currently unspent.
func (self *TrOutput)GetSpentIn()string{
  return self.spent
}

// Store the hash of a block as the block in which this
// output is spent.
func (self *TrOutput)SetSpentIn(block string){
  self.spent = block
}
