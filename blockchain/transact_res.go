/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-22
 * @Project: Proof of Evolution
 * @Filename: transact_res.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-22
 * @Copyright: 2020
 */

package blockchain

import(
  "time"
  "github.com/D33pBlue/poe/utils"
)

type ResTransaction struct{
  Timestamp time.Time
  Inputs []TrInput
  Output TrOutput // the possible remainder
  Creator utils.Addr
  JobBlock string // the hash of the block with the JobTransaction
  JobTrans string // the hash of the JobTransaction
  Evaluation float64
  HashSol string // the hash of solution: [individual,evaluation,creator]
  Hash string
  Signature string
}

func (self *ResTransaction)Check(block *Block,trChanges *map[string]string)bool{
  return false // TODO: implement later
}

func (self *ResTransaction)GetHash()string{
  return "" // TODO: implement later
}

func (self *ResTransaction)GetHashCached()string{
  return self.Hash
}

func (self *ResTransaction)GetCreator()utils.Addr{
  return self.Creator
}

func (self *ResTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

func (self *ResTransaction)GetType()string{
  return TrRes
}

func (self *ResTransaction)GetOutputAt(int)*TrOutput{
  return &self.Output
}

func (self *ResTransaction)Serialize()[]byte{
  return nil // TODO: implement later
}

func MarshalResTransaction([]byte)*ResTransaction{
  return nil // TODO: implement later
}
